package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-clean/platform/config"
	"github.com/go-clean/platform/logger"
	"github.com/redis/go-redis/v9"
)

// NewClient creates a new Redis client
func NewClient(cfg config.RedisConfig, log logger.Logger) (*redis.Client, error) {
	log.Info().Str("host", cfg.Host).Int("port", cfg.Port).Int("db", cfg.DB).Msg("Initializing Redis connection")

	// Create Redis client options
	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	}

	// Create Redis client
	log.Debug().Int("pool_size", cfg.PoolSize).Int("min_idle_conns", cfg.MinIdleConns).Msg("Creating Redis client")
	client := redis.NewClient(opts)

	// Test the connection
	log.Debug().Msg("Testing Redis connection")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Redis ping failed")
		client.Close()
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info().Msg("Redis client created successfully")
	return client, nil
}

// Close gracefully closes the Redis client
func Close(client *redis.Client, log logger.Logger) error {
	if client != nil {
		log.Info().Msg("Closing Redis client")
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close Redis client")
		} else {
			log.Debug().Msg("Redis client closed successfully")
		}
		return err
	}
	return nil
}
