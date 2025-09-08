//go:build wireinject
// +build wireinject

package ratelimit

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	
	"github.com/go-clean/internal/ratelimit/application/command"
	"github.com/go-clean/internal/ratelimit/infrastructure"
	"github.com/go-clean/internal/ratelimit/ports"
	"github.com/go-clean/internal/ratelimit/presentation/http"
	"github.com/go-clean/platform/logger"
)

// ProviderSet is the Wire provider set for the rate-limit module (Redis only)
var ProviderSet = wire.NewSet(
	// Infrastructure providers
	infrastructure.NewRedisRateLimitRepository,
	wire.Bind(new(ports.RateLimitRepository), new(*infrastructure.RedisRateLimitRepository)),
	
	// Application providers
	command.NewCheckRateLimitCommandHandler,
	command.NewCheckRateLimitWithDetailCommandHandler,
	
	// Presentation providers
	http.NewRateLimitHandler,
)

// HybridProviderSet is the Wire provider set for the rate-limit module with hybrid caching
var HybridProviderSet = wire.NewSet(
	// Infrastructure providers
	infrastructure.NewRedisRateLimitRepository,
	infrastructure.NewHybridRateLimitRepository,
	wire.Bind(new(ports.RateLimitRepository), new(*infrastructure.HybridRateLimitRepository)),
	
	// Application providers
	command.NewCheckRateLimitCommandHandler,
	command.NewCheckRateLimitWithDetailCommandHandler,
	
	// Presentation providers
	http.NewRateLimitHandler,
)

// NewRateLimitModule creates a new rate-limit module with all dependencies wired
func NewRateLimitModule(
	logger logger.Logger,
	redisClient *redis.Client,
) (*http.RateLimitHandler, error) {
	wire.Build(ProviderSet)
	return nil, nil
}