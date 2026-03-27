package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	// CSPNonceKey is the context key for storing the CSP nonce
	CSPNonceKey = "csp_nonce"
	// NonceTemplate is the placeholder in CSP policy for nonce
	NonceTemplate = "__CSP_NONCE__"
	// CloudflareInsightsDomain is the domain for Cloudflare Web Analytics
	CloudflareInsightsDomain = "https://static.cloudflareinsights.com"
)

// GenerateNonce generates a cryptographically secure random nonce.
// 返回 error 以确保调用方在 crypto/rand 失败时能正确降级。
func GenerateNonce() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate CSP nonce: %w", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GetNonceFromContext retrieves the CSP nonce from gin context
func GetNonceFromContext(c *gin.Context) string {
	if nonce, exists := c.Get(CSPNonceKey); exists {
		if s, ok := nonce.(string); ok {
			return s
		}
	}
	return ""
}

// SecurityHeaders sets baseline security headers for all responses.
// getFrameSrcOrigins is an optional function that returns extra origins to inject into frame-src;
// pass nil to disable dynamic frame-src injection.
func SecurityHeaders(cfg config.CSPConfig, getFrameSrcOrigins func() []string, hstsMaxAge int64) gin.HandlerFunc {
	policy := strings.TrimSpace(cfg.Policy)
	if policy == "" {
		policy = config.DefaultCSPPolicy
	}

	// Enhance policy with required directives (nonce placeholder and Cloudflare Insights)
	policy = enhanceCSPPolicy(policy)

	return func(c *gin.Context) {
		finalPolicy := policy
		if getFrameSrcOrigins != nil {
			for _, origin := range getFrameSrcOrigins() {
				if origin != "" {
					finalPolicy = addToDirective(finalPolicy, "frame-src", origin)
				}
			}
		}

		// ── Core Security Headers ──────────────────────────────────
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Prevent clickjacking (API routes skip frame-ancestors — already enforced by CSP)
		c.Header("X-Frame-Options", "DENY")
		// Referrer policy for cross-origin requests
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Control information in browser address bar
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// ── Transport Security (HSTS) ─────────────────────────────
		// Only send HSTS on HTTPS requests; skip upgrade-insecure-requests
		// so that the server doesn't redirect API clients (they expect raw responses).
		scheme := c.GetHeader("X-Forwarded-Proto")
		if scheme == "" {
			// Infer from the request TLS state
			if c.Request.TLS != nil {
				scheme = "https"
			} else {
				scheme = "http"
			}
		}
		if scheme == "https" && hstsMaxAge > 0 {
			// includeSubDomains and preload are not set by default to avoid accidental
			// subdomain breakage; operators can extend the header via config if needed.
			c.Header("Strict-Transport-Security",
				fmt.Sprintf("max-age=%d; includeSubDomains", hstsMaxAge))
		}

		// ── Browser Feature / Permissions Policy ─────────────────
		// Restrict powerful browser APIs to same-origin only.
		// Omit deprecated/removed tokens (e.g. document-domain, speaker) that Chrome warns on.
		c.Header("Permissions-Policy",
			"accelerometer=(), camera=(), geolocation=(), gyroscope=(), "+
				"magnetometer=(), microphone=(), payment=(self), sync-xhr=(), usb=()")

		apiPath := isAPIRoutePath(c)

		// ── Cross-Origin Isolation (COEP/COOP) ───────────────────
		// Apply only to HTML/SPA responses. JSON API responses under /api/v1 should not
		// carry COEP/COOP (unnecessary and can confuse clients / tooling).
		if !apiPath {
			c.Header("Cross-Origin-Embedder-Policy", "require-corp")
			c.Header("Cross-Origin-Opener-Policy", "same-origin")
		}

		// ── API Routes: skip CSP/nonce processing ─────────────────
		if apiPath {
			c.Next()
			return
		}

		// ── CSP for frontend routes ───────────────────────────────
		if cfg.Enabled {
			// Generate nonce for this request
			nonce, err := GenerateNonce()
			if err != nil {
				// crypto/rand 失败时降级为无 nonce 的 CSP 策略
				log.Printf("[SecurityHeaders] %v — 降级为无 nonce 的 CSP", err)
				c.Header("Content-Security-Policy", strings.ReplaceAll(finalPolicy, NonceTemplate, "'unsafe-inline'"))
			} else {
				c.Set(CSPNonceKey, nonce)
				c.Header("Content-Security-Policy", strings.ReplaceAll(finalPolicy, NonceTemplate, "'nonce-"+nonce+"'"))
			}
		}
		c.Next()
	}
}

func isAPIRoutePath(c *gin.Context) bool {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return false
	}
	path := c.Request.URL.Path
	return strings.HasPrefix(path, "/api/v1/") ||
		strings.HasPrefix(path, "/v1/") ||
		strings.HasPrefix(path, "/v1beta/") ||
		strings.HasPrefix(path, "/antigravity/") ||
		strings.HasPrefix(path, "/sora/") ||
		strings.HasPrefix(path, "/responses")
}

// enhanceCSPPolicy ensures the CSP policy includes nonce support and Cloudflare Insights domain.
// This allows the application to work correctly even if the config file has an older CSP policy.
func enhanceCSPPolicy(policy string) string {
	// Add nonce placeholder to script-src if not present
	if !strings.Contains(policy, NonceTemplate) && !strings.Contains(policy, "'nonce-") {
		policy = addToDirective(policy, "script-src", NonceTemplate)
	}

	// Add Cloudflare Insights domain to script-src if not present
	if !strings.Contains(policy, CloudflareInsightsDomain) {
		policy = addToDirective(policy, "script-src", CloudflareInsightsDomain)
	}

	return policy
}

// addToDirective adds a value to a specific CSP directive.
// If the directive doesn't exist, it will be added after default-src.
func addToDirective(policy, directive, value string) string {
	// Find the directive in the policy
	directivePrefix := directive + " "
	idx := strings.Index(policy, directivePrefix)

	if idx == -1 {
		// Directive not found, add it after default-src or at the beginning
		defaultSrcIdx := strings.Index(policy, "default-src ")
		if defaultSrcIdx != -1 {
			// Find the end of default-src directive (next semicolon)
			endIdx := strings.Index(policy[defaultSrcIdx:], ";")
			if endIdx != -1 {
				insertPos := defaultSrcIdx + endIdx + 1
				// Insert new directive after default-src
				return policy[:insertPos] + " " + directive + " 'self' " + value + ";" + policy[insertPos:]
			}
		}
		// Fallback: prepend the directive
		return directive + " 'self' " + value + "; " + policy
	}

	// Find the end of this directive (next semicolon or end of string)
	endIdx := strings.Index(policy[idx:], ";")

	if endIdx == -1 {
		// No semicolon found, directive goes to end of string
		return policy + " " + value
	}

	// Insert value before the semicolon
	insertPos := idx + endIdx
	return policy[:insertPos] + " " + value + policy[insertPos:]
}
