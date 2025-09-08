package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-clean/platform/logger"
)

// RedisRateLimitRepository implements the RateLimitRepository interface using Redis
type RedisRateLimitRepository struct {
	logger      logger.Logger
	redisClient *redis.Client
	windowSize  time.Duration
}

// NewRedisRateLimitRepository creates a new Redis-based rate limit repository
func NewRedisRateLimitRepository(
	logger logger.Logger,
	redisClient *redis.Client,
) *RedisRateLimitRepository {
	return &RedisRateLimitRepository{
		logger:      logger,
		redisClient: redisClient,
		windowSize:  time.Minute, // Default 1-minute window
	}
}

// RateLimit checks if a user is allowed to make a request based on the rate limit
func (r *RedisRateLimitRepository) RateLimit(userId string, limit int) bool {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:%s", userId)

	r.logger.Debug().Str("user_id", userId).Int("limit", limit).Str("key", key).Msg("Checking rate limit")

	// Use Redis pipeline for atomic operations
	pipe := r.redisClient.Pipeline()

	// Set expiration if this is the first request
	_ = pipe.SetNX(ctx, key, 0, r.windowSize)
	// Increment the counter
	incrCmd := pipe.Incr(ctx, key)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error().Str("user_id", userId).Err(err).Msg("Failed to execute Redis pipeline")
		return false // Fail closed - deny request on error
	}

	// Get the current count
	currentCount := incrCmd.Val()

	r.logger.Debug().Str("user_id", userId).Int64("current_count", currentCount).Int("limit", limit).Bool("allowed", currentCount <= int64(limit)).Msg("Rate limit check result")

	return currentCount <= int64(limit)
}

// RateLimitWithDetail checks if a user is allowed to make a request and returns detailed information
func (r *RedisRateLimitRepository) RateLimitWithDetail(userId string, limit int) (int, time.Duration, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:%s", userId)

	r.logger.Debug().Str("user_id", userId).Int("limit", limit).Str("key", key).Msg("Checking rate limit with detail")

	// Use Redis pipeline for atomic operations
	pipe := r.redisClient.Pipeline()

	_ = pipe.SetNX(ctx, key, 0, r.windowSize)
	// Increment the counter
	incrCmd := pipe.Incr(ctx, key)
	// Get TTL for reset time calculation
	ttlCmd := pipe.TTL(ctx, key)

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error().Str("user_id", userId).Err(err).Msg("Failed to execute Redis pipeline")
		return 0, 0, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Get the current count and TTL
	currentCount := incrCmd.Val()
	ttl := ttlCmd.Val()

	// Calculate remaining requests
	remaining := limit - int(currentCount)
	if remaining < 0 {
		remaining = 0
	}

	r.logger.Debug().Str("user_id", userId).Int64("current_count", currentCount).Int("limit", limit).Int("remaining", remaining).Dur("ttl", ttl).Msg("Rate limit check with detail result")

	return remaining, ttl, nil
}
