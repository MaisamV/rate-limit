package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-clean/platform/config"
	"github.com/go-clean/platform/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new PostgreSQL connection pool
func NewConnection(cfg config.DatabaseConfig, log logger.Logger) (*pgxpool.Pool, error) {
	log.Info().Str("host", cfg.Host).Int("port", cfg.Port).Str("database", cfg.DBName).Msg("Initializing database connection")

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse database configuration")
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	log.Debug().Int("max_conns", cfg.MaxOpenConns).Int("max_idle_conns", cfg.MaxIdleConns).Msg("Creating database connection pool")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database connection pool")
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// Test the connection
	log.Debug().Msg("Testing database connection")
	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Database ping failed")
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection pool created successfully")
	return pool, nil
}

// Close gracefully closes the database connection pool
func Close(pool *pgxpool.Pool, log logger.Logger) {
	if pool != nil {
		log.Info().Msg("Closing database connection pool")
		pool.Close()
		log.Debug().Msg("Database connection pool closed")
	}
}
