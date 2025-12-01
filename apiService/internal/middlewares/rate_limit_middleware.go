package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	// Requests per window
	Limit int
	// Window duration
	Window time.Duration
	// Key prefix for Redis
	KeyPrefix string
}

// DefaultAPIRateLimitConfig returns default config for API endpoints
func DefaultAPIRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Limit:     200,         // 200 requests
		Window:    time.Minute, // per minute
		KeyPrefix: "ratelimit:user:",
	}
}

// StrictRateLimitConfig returns strict config for sensitive endpoints
func StrictRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Limit:     20,          // 20 requests
		Window:    time.Minute, // per minute
		KeyPrefix: "ratelimit:user:strict:",
	}
}

// RateLimitMiddleware creates a per-user rate limiting middleware using Redis
// Uses sliding window algorithm for accurate rate limiting
func RateLimitMiddleware(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем OPTIONS запросы (preflight)
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Skip if Redis is not available
		if redisClient == nil {
			c.Next()
			return
		}

		// Get userID from context (set by JWT middleware)
		userIDValue, exists := c.Get("userID")
		if !exists {
			// No user ID - skip per-user limiting (Nginx handles per-IP)
			c.Next()
			return
		}

		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Create rate limit key
		key := fmt.Sprintf("%s%s", config.KeyPrefix, userID.String())

		// Check and increment counter using Redis
		allowed, remaining, resetAt, err := checkRateLimit(ctx, redisClient, key, config)
		if err != nil {
			// On Redis error, allow request (fail open)
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))

		if !allowed {
			c.Header("Retry-After", strconv.FormatInt(resetAt-time.Now().Unix(), 10))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Too many requests. Limit: %d requests per %s", config.Limit, config.Window),
			})
			return
		}

		c.Next()
	}
}

// checkRateLimit implements sliding window rate limiting
func checkRateLimit(ctx context.Context, client *redis.Client, key string, config RateLimitConfig) (allowed bool, remaining int, resetAt int64, err error) {
	now := time.Now()
	windowStart := now.Add(-config.Window)
	resetAt = now.Add(config.Window).Unix()

	// Use Redis transaction for atomic operations
	pipe := client.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.UnixNano(), 10))

	// Count current requests in window
	countCmd := pipe.ZCard(ctx, key)

	// Execute pipeline
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return true, config.Limit, resetAt, err
	}

	currentCount := int(countCmd.Val())

	// Check if limit exceeded
	if currentCount >= config.Limit {
		remaining = 0
		return false, remaining, resetAt, nil
	}

	// Add new request to sorted set
	member := fmt.Sprintf("%d:%s", now.UnixNano(), uuid.New().String()[:8])
	err = client.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: member,
	}).Err()
	if err != nil {
		return true, config.Limit - currentCount, resetAt, err
	}

	// Set expiration on key
	client.Expire(ctx, key, config.Window+time.Second)

	remaining = config.Limit - currentCount - 1
	if remaining < 0 {
		remaining = 0
	}

	return true, remaining, resetAt, nil
}

// RateLimitByEndpoint creates rate limiter with endpoint-specific keys
func RateLimitByEndpoint(redisClient *redis.Client, config RateLimitConfig, endpoint string) gin.HandlerFunc {
	endpointConfig := RateLimitConfig{
		Limit:     config.Limit,
		Window:    config.Window,
		KeyPrefix: fmt.Sprintf("%s%s:", config.KeyPrefix, endpoint),
	}
	return RateLimitMiddleware(redisClient, endpointConfig)
}
