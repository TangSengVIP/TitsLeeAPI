package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// APIGatewayRateLimiter provides multi-dimensional rate limiting for API gateway routes.
// It limits on two independent axes simultaneously:
//   - By (API Key or IP): per-key/IP request count within a sliding window
//   - By endpoint pattern: per-endpoint global burst protection
//
// Both limits must pass for a request to proceed.
type APIGatewayRateLimiter struct {
	redis  *redis.Client
	prefix string
}

// NewAPIGatewayRateLimiter creates a gateway rate limiter backed by Redis.
func NewAPIGatewayRateLimiter(redisClient *redis.Client) *APIGatewayRateLimiter {
	return &APIGatewayRateLimiter{
		redis:  redisClient,
		prefix: "gwr:",
	}
}

// GatewayRateLimitConfig controls the per-axis limits.
type GatewayRateLimitConfig struct {
	// PerKeyLimit: max requests per key/IP in the window.
	// Set to 0 to disable this axis.
	PerKeyLimit int
	// PerEndpointLimit: max requests to this specific endpoint globally in the window.
	// Protects against burst traffic to a single route.
	// Set to 0 to disable.
	PerEndpointLimit int
	// Window: time window for the limit.
	Window time.Duration
	// FailureMode: what to do when Redis is unavailable.
	// Recommended: RateLimitFailClose for sensitive endpoints, RateLimitFailOpen for bulk reads.
	FailureMode RateLimitFailureMode
}

// gatewayRateLimitLua implements atomic multi-key rate limiting with a single Lua script.
// Returns: {per_key_count, per_endpoint_count, per_key_repaired, per_endpoint_repaired}
var gatewayRateLimitLua = redis.NewScript(`
local function set_ttl_if_needed(key, window_ms)
  local ttl = redis.call('PTTL', key)
  if ttl == -1 then
    redis.call('PEXPIRE', key, window_ms)
    return 1
  end
  return 0
end

-- Per-key counter
local key_count = redis.call('INCR', KEYS[1])
set_ttl_if_needed(KEYS[1], ARGV[1])
local key_repaired = 0
if key_count == 1 then
  key_repaired = 1
end

-- Per-endpoint counter
local ep_count = redis.call('INCR', KEYS[2])
set_ttl_if_needed(KEYS[2], ARGV[1])
local ep_repaired = 0
if ep_count == 1 then
  ep_repaired = 1
end

return {key_count, ep_count, key_repaired, ep_repaired}
`)

// Limit returns a middleware that enforces per-key/IP and per-endpoint rate limits.
// clientIP is extracted using the ip package's trusted IP resolution.
func (r *APIGatewayRateLimiter) Limit(cfg GatewayRateLimitConfig) gin.HandlerFunc {
	if cfg.PerKeyLimit <= 0 && cfg.PerEndpointLimit <= 0 {
		// Both axes disabled — return a no-op
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		windowMillis := cfg.Window.Milliseconds()
		if windowMillis < 1 {
			windowMillis = 60000
		}

		// Determine the key identifier
		var keyID string
		if apiKey, ok := servermiddleware.GetAPIKeyFromContext(c); ok && apiKey != nil {
			// Authenticated: rate limit by API key (most granular)
			keyID = fmt.Sprintf("key:%d", apiKey.ID)
		} else {
			// Unauthenticated: rate limit by client IP
			clientIP := ip.GetTrustedClientIP(c)
			keyID = fmt.Sprintf("ip:%s", clientIP)
		}

		// Determine the endpoint identifier
		method := c.Request.Method
		path := normalizeEndpointPath(c.Request.URL.Path)
		endpointID := fmt.Sprintf("ep:%s:%s", method, path)

		var keyCount, epCount int64
		var err error

		if cfg.PerKeyLimit > 0 && cfg.PerEndpointLimit > 0 {
			// Both axes enabled — use atomic Lua script
			key1 := r.prefix + keyID
			key2 := r.prefix + endpointID
			vals, luaErr := gatewayRateLimitLua.Run(ctx, r.redis, []string{key1, key2}, windowMillis).Slice()
			if luaErr != nil {
				log.Printf("[GatewayRateLimit] lua script error: %v", luaErr)
				err = luaErr
			} else if len(vals) >= 2 {
				keyCount, _ = parseLuaInt64(vals[0])
				epCount, _ = parseLuaInt64(vals[1])
			}
		} else if cfg.PerKeyLimit > 0 {
			// Only per-key limit
			keyCount, _, err = rateLimitRun(ctx, r.redis, r.prefix+keyID, windowMillis)
		} else {
			// Only per-endpoint limit
			epCount, _, err = rateLimitRun(ctx, r.redis, r.prefix+endpointID, windowMillis)
		}

		if err != nil {
			log.Printf("[GatewayRateLimit] redis error: key=%s err=%v mode=%s",
				keyID, err, failureModeLabel(cfg.FailureMode))
			if cfg.FailureMode == RateLimitFailClose {
				abortGatewayLimit(c)
				return
			}
			// Fail open — let the request through
			c.Next()
			return
		}

		// Check per-key limit
		if cfg.PerKeyLimit > 0 && keyCount > int64(cfg.PerKeyLimit) {
			log.Printf("[GatewayRateLimit] per-key limit exceeded: key=%s count=%d limit=%d",
				keyID, keyCount, cfg.PerKeyLimit)
			abortGatewayLimit(c)
			return
		}

		// Check per-endpoint limit
		if cfg.PerEndpointLimit > 0 && epCount > int64(cfg.PerEndpointLimit) {
			log.Printf("[GatewayRateLimit] per-endpoint limit exceeded: endpoint=%s:%s count=%d limit=%d",
				method, path, epCount, cfg.PerEndpointLimit)
			abortGatewayLimit(c)
			return
		}

		// Add rate limit headers so clients can adapt
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.PerKeyLimit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(max(0, int64(cfg.PerKeyLimit)-keyCount), 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(cfg.Window).Unix(), 10))

		c.Next()
	}
}

// normalizeEndpointPath normalizes a path for use as a rate limit key.
// Strips variable segments (e.g. /users/123 -> /users/{id}) to avoid
// creating a new bucket for each unique resource.
func normalizeEndpointPath(path string) string {
	// Strip trailing slashes for consistency
	path = strings.TrimSuffix(path, "/")
	// Collapse multiple slashes
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	return path
}

func abortGatewayLimit(c *gin.Context) {
	c.Header("Retry-After", "60")
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"type":    "error",
		"error": gin.H{
			"type":    "rate_limit_error",
			"message": "Rate limit exceeded. Please slow down and retry later.",
		},
	})
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func parseLuaInt64(v any) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected lua type %T", v)
	}
}
