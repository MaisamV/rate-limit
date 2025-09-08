package domain

import (
	"time"
)

// RateLimit represents a rate limit entity
type RateLimit struct {
	UserID    string
	Limit     int
	Remaining int
	ResetTime time.Time
	Window    time.Duration
}

// NewRateLimit creates a new rate limit instance
func NewRateLimit(userID string, limit int, window time.Duration) *RateLimit {
	return &RateLimit{
		UserID:    userID,
		Limit:     limit,
		Remaining: limit,
		ResetTime: time.Now().Add(window),
		Window:    window,
	}
}

// IsAllowed checks if the request is allowed based on current rate limit
func (rl *RateLimit) IsAllowed() bool {
	if time.Now().After(rl.ResetTime) {
		// Reset the rate limit window
		rl.Remaining = rl.Limit
		rl.ResetTime = time.Now().Add(rl.Window)
	}

	if rl.Remaining > 0 {
		rl.Remaining--
		return true
	}

	return false
}

// Reset resets the rate limit to its initial state
func (rl *RateLimit) Reset() {
	rl.Remaining = rl.Limit
	rl.ResetTime = time.Now().Add(rl.Window)
}