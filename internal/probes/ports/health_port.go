package ports

import (
	"context"
	"time"
)

// DatabaseChecker defines the interface for checking database connectivity
type DatabaseChecker interface {
	CheckDatabase(ctx context.Context) (bool, time.Duration, error)
}

// RedisChecker defines the interface for checking Redis connectivity
type RedisChecker interface {
	CheckRedis(ctx context.Context) (bool, time.Duration, error)
}
