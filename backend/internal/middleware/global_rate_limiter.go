package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// GlobalRateLimiter provides IP-based rate limiting across all routes.
// This is a defense-in-depth measure against:
//   - Distributed brute-force attacks on auth endpoints
//   - Web scanning and crawling
//   - General DoS from a single source IP
type GlobalRateLimiter struct {
	redis  *redis.Client
	prefix string
}

// NewGlobalRateLimiter creates a global IP-based rate limiter.
func NewGlobalRateLimiter(redisClient *redis.Client) *GlobalRateLimiter {
	return &GlobalRateLimiter{
		redis:  redisClient,
		prefix: "grl:",
	}
}

// Limit returns a middleware that enforces a global per-IP rate limit.
// This should be applied early in the middleware chain, before authentication,
// to catch unauthenticated abuse.
func (r *GlobalRateLimiter) Limit(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get trusted client IP (uses X-Forwarded-For when trusted proxies are configured)
		clientIP := c.ClientIP()
		key := r.prefix + clientIP

		ctx := c.Request.Context()
		windowMillis := window.Milliseconds()
		if windowMillis < 1 {
			windowMillis = 60000
		}

		count, _, err := rateLimitRun(ctx, r.redis, key, windowMillis)
		if err != nil {
			// On Redis failure, fail open to avoid blocking legitimate traffic.
			// Individual endpoints have their own fail-close rate limiters.
			c.Next()
			return
		}

		if count > int64(requests) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests from this IP, please try again later",
			})
			return
		}

		c.Next()
	}
}

// buildGlobalLimitMw creates a Gin middleware that enforces a global per-IP
// rate limit of 200 req/min per IP (fail-open for availability).
func BuildGlobalLimitMw(redisClient *redis.Client) gin.HandlerFunc {
	lr := NewGlobalRateLimiter(redisClient)
	return lr.Limit(200, time.Minute)
}

// AbuseDetectionLimiter is a specialized rate limiter for auth-related endpoints.
// It uses tighter limits and fail-close behavior to prevent credential stuffing.
// Apply this to the auth group before the per-endpoint rate limiters.
type AbuseDetectionLimiter struct {
	redis  *redis.Client
	prefix string
}

// NewAbuseDetectionLimiter creates an abuse detection rate limiter.
func NewAbuseDetectionLimiter(redisClient *redis.Client) *AbuseDetectionLimiter {
	return &AbuseDetectionLimiter{
		redis:  redisClient,
		prefix: "abuse:",
	}
}

// FailCloseLimit returns a middleware that enforces a fail-close per-IP rate limit.
// If Redis is unavailable, requests are BLOCKED (fail-close), making it suitable
// for protecting high-value unauthenticated endpoints.
func (r *AbuseDetectionLimiter) FailCloseLimit(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := r.prefix + clientIP

		ctx := c.Request.Context()
		windowMillis := window.Milliseconds()
		if windowMillis < 1 {
			windowMillis = 60000
		}

		count, _, err := rateLimitRun(ctx, r.redis, key, windowMillis)
		if err != nil {
			// Fail-close: block if Redis is unavailable
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":    "SERVICE_UNAVAILABLE",
				"message": "Service temporarily unavailable",
			})
			return
		}

		if count > int64(requests) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    "ABUSE_DETECTED",
				"message": "Too many requests, please try again later",
			})
			return
		}

		c.Next()
	}
}
