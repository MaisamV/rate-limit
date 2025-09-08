package ports

import "time"

// RateLimitRepository defines the interface for rate limit data access
type RateLimitRepository interface {
	// RateLimit checks if a user is allowed to make a request based on the rate limit
	// Returns true if the request is allowed, false otherwise
	RateLimit(userId string, limit int) bool
	
	// RateLimitWithDetail checks if a user is allowed to make a request and returns detailed information
	// Returns current count, time until reset, and error if any
	RateLimitWithDetail(userId string, limit int) (int, time.Duration, error)
}