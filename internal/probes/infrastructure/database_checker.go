package infrastructure

import (
	"context"
	"time"

	"github.com/go-clean/platform/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseChecker implements the DatabaseChecker port
type DatabaseChecker struct {
	logger logger.Logger
	db     *pgxpool.Pool
}

// NewDatabaseChecker creates a new database checker
func NewDatabaseChecker(logger logger.Logger, db *pgxpool.Pool) *DatabaseChecker {
	return &DatabaseChecker{
		logger: logger,
		db:     db,
	}
}

// CheckDatabase checks the database connectivity and response time
func (dc *DatabaseChecker) CheckDatabase(ctx context.Context) (bool, time.Duration, error) {
	dc.logger.Debug().Msg("Starting database connectivity check")
	start := time.Now()

	// Create a context with timeout for the health check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Simple ping to check database connectivity
	err := dc.db.Ping(checkCtx)
	duration := time.Since(start)

	if err != nil {
		dc.logger.Error().Err(err).Int64("duration_ms", duration.Milliseconds()).Msg("Database ping failed")
		return false, duration, err
	}

	dc.logger.Debug().Int64("duration_ms", duration.Milliseconds()).Msg("Database ping successful")
	return true, duration, nil
}
