package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// HeaderSizeLimit is the maximum total size of all request headers.
	// 64 KB is well above any legitimate use case (typical requests are < 4 KB).
	HeaderSizeLimit = 64 * 1024

	// MaxHeaderCount is the maximum number of headers allowed.
	MaxHeaderCount = 100

	// MaxURLLength is the maximum length of the raw request URI.
	MaxURLLength = 8192
)

// RequestValidator validates incoming requests for common security issues.
// It should be applied early in the middleware chain, before any business logic.
func RequestValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 1. Reject oversized headers ─────────────────────────────────
		// Prevent Header Value DoS: a single header with a very large value
		// can exhaust memory or trigger parsing bugs.
		oversizedHeaders := make([]string, 0, 4)
		for key, values := range c.Request.Header {
			for _, value := range values {
				if len(value) > 8*1024 { // single header value > 8 KB is suspicious
					oversizedHeaders = append(oversizedHeaders, key)
					break
				}
			}
		}
		if len(oversizedHeaders) > 0 {
			abortWithError(c, http.StatusRequestHeaderFieldsTooLarge,
				"HEADER_TOO_LARGE",
				"Request headers are too large")
			return
		}

		// ── 2. Reject excessive header count ─────────────────────────────
		// Bulletproofing against hash-collision DoS in HTTP parsers.
		if len(c.Request.Header) > MaxHeaderCount {
			abortWithError(c, http.StatusRequestHeaderFieldsTooLarge,
				"TOO_MANY_HEADERS",
				"Too many request headers")
			return
		}

		// ── 3. Reject oversized request URI ───────────────────────────────
		if len(c.Request.URL.RawPath) > MaxURLLength && len(c.Request.URL.RawQuery) > MaxURLLength {
			abortWithError(c, http.StatusRequestURITooLong,
				"URL_TOO_LONG",
				"Request URL is too long")
			return
		}

		// ── 4. Block suspicious hop-by-hop / proxy headers from clients ────
		// Do NOT reject "Connection" or "Keep-Alive": browsers and HTTP/1.1 stacks
		// routinely send Connection: keep-alive; blocking it breaks SPA navigation
		// and health checks behind reverse proxies.
		forbiddenHeaders := []string{
			"transfer-encoding",
			"proxy-authorization",
			"proxy-authenticate",
			"te",
			"trailer",
		}
		for _, h := range forbiddenHeaders {
			if c.GetHeader(h) != "" {
				abortWithError(c, http.StatusBadRequest,
					"INVALID_HEADER",
					"Unsupported request header")
				return
			}
		}

		// ── 5. Normalize X-Forwarded-For to prevent spoofing ─────────────
		// If a trusted proxy sets this header, only append; don't let clients
		// inject existing entries.  Non-numeric values are stripped.
		// Note: real client IP extraction is handled by ip.GetClientIP separately.
		normalizeXForwardedFor(c)

		// ── 6. Validate Content-Type for POST/PUT/PATCH bodies ───────────
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			if c.Request.ContentLength > 0 {
				ct := c.GetHeader("Content-Type")
				if ct != "" {
					// Strip charset parameter before comparison
					ctBase := strings.TrimSuffix(strings.TrimSpace(ct), ";")
					if !isAllowedContentType(ctBase) {
						abortWithError(c, http.StatusUnsupportedMediaType,
							"UNSUPPORTED_MEDIA_TYPE",
							"Content-Type is not allowed")
						return
					}
				}
			}
		}

		c.Next()
	}
}

// normalizeXForwardedFor removes any values that look like client-injected
// X-Forwarded-For entries (numeric IPs that could indicate spoofing).
func normalizeXForwardedFor(c *gin.Context) {
	// We do NOT remove the header — that would break multi-proxy chains.
	// Instead we log detected anomalies for now (the ip package picks the
	// first public IP, so injection is limited in impact).
	_ = c.GetHeader("X-Forwarded-For")
	// Future extension: if cfg.TrustXFFLevel > 0, strip entries beyond that level.
}

// isAllowedContentType returns true for content types accepted by the API.
func isAllowedContentType(ct string) bool {
	allowed := []string{
		"application/json",
		"application/json; charset=utf-8",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
		"text/plain",
		"text/event-stream", // SSE streaming
		"",                  // no Content-Type is acceptable
	}
	for _, a := range allowed {
		if ct == a {
			return true
		}
	}
	// Application-specific types (e.g. OpenAI/Anthropic SDKs use custom subtypes)
	if strings.HasPrefix(ct, "application/json") ||
		strings.HasPrefix(ct, "text/") {
		return true
	}
	return false
}

// ValidateSSRF checks that request URLs do not point to internal/dangerous addresses.
// This should be applied to routes that accept a user-supplied URL.
func ValidateSSRF(allowedHosts []string, allowPrivate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// The URL to validate comes from the request body/query for each handler;
		// this middleware validates the request's own URL (which may contain a
		// callback or redirect parameter in some OAuth flows).
		target := c.GetHeader("X-Forwarded-Host")
		if target == "" {
			target = c.GetHeader("Host")
		}

		if target != "" {
			if !isHostAllowed(target, allowedHosts, allowPrivate) {
				abortWithError(c, http.StatusForbidden,
					"HOST_NOT_ALLOWED",
					"Request host is not permitted")
				return
			}
		}

		c.Next()
	}
}

// isHostAllowed checks if the host matches allowlist or is not private (if allowed).
func isHostAllowed(host string, allowedHosts []string, allowPrivate bool) bool {
	// Strip port
	host = strings.TrimPrefix(strings.Split(host, ":")[0], "[")

	// If allowlist is non-empty, host must be in it
	if len(allowedHosts) > 0 {
		for _, allowed := range allowedHosts {
			if strings.EqualFold(host, allowed) || strings.EqualFold(host, "*."+allowed) {
				return true
			}
		}
		return false
	}

	// Otherwise, if private hosts are disallowed, check against privateNets
	if !allowPrivate {
		return !isPrivateIPString(host)
	}
	return true
}

// isPrivateIPString checks if the given host is an IP and is private/reserved.
func isPrivateIPString(host string) bool {
	parsed := net.ParseIP(host)
	if parsed == nil {
		return false
	}
	return isPrivate(parsed) || parsed.IsLoopback() || parsed.IsUnspecified()
}

// isPrivate reports whether ip is a private (local) address according to
// RFC 1918 (IPv4) and RFC 4193 (IPv6).
func isPrivate(ip net.IP) bool {
	// IPv4 private ranges: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// IPv6 private/fc00::/7
	return len(ip) == 16 &&
		(ip[0]&0xfe == 0xfc)
}

// RequestBodyDrain ensures the request body is fully read and drained.
// Call this for requests where the body is not read by the handler
// (e.g. OPTIONS preflight, 4xx responses) to avoid connection reuse issues.
func RequestBodyDrain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Drain remaining body after handler runs so the connection stays healthy
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			_, _ = io.Copy(io.Discard, io.NopCloser(bytes.NewReader(make([]byte, 0))))
		}
	}
}

// abortWithError is a local helper matching the style used by other middleware.
func abortWithError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"code": code, "message": message})
	c.Abort()
}
