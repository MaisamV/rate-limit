package infrastructure

import (
	"context"
	"time"

	"github.com/go-clean/platform/logger"
	"github.com/redis/go-redis/v9"
)

// RedisChecker implements the RedisChecker port
type RedisChecker struct {
	logger logger.Logger
	client *redis.Client
}

// NewRedisChecker creates a new Redis checker
func NewRedisChecker(logger logger.Logger, client *redis.Client) *RedisChecker {
	return &RedisChecker{
		logger: logger,
		client: client,
	}
}

// CheckRedis checks the Redis connectivity and response time
func (rc *RedisChecker) CheckRedis(ctx context.Context) (bool, time.Duration, error) {
	rc.logger.Debug().Msg("Starting Redis connectivity check")
	start := time.Now()

	// Create a context with timeout for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple ping to check Redis connectivity
	result := rc.client.Ping(checkCtx)
	duration := time.Since(start)

	if err := result.Err(); err != nil {
		rc.logger.Error().Err(err).Int64("duration_ms", duration.Milliseconds()).Msg("Redis ping failed")
		return false, duration, err
	}

	rc.logger.Debug().Int64("duration_ms", duration.Milliseconds()).Msg("Redis ping successful")
	return true, duration, nil
}
