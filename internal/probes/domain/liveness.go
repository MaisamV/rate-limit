package domain

import "time"

// LivenessStatus represents the liveness status
type LivenessStatus string

const (
	LivenessStatusAlive LivenessStatus = "alive"
	LivenessStatusDead  LivenessStatus = "dead"
)

// LivenessResponse represents the liveness check response
type LivenessResponse struct {
	Status        LivenessStatus `json:"status"`
	UptimeSeconds int64          `json:"uptime_seconds"`
	Timestamp     time.Time      `json:"timestamp"`
}

// NewLivenessResponse creates a new liveness response
func NewLivenessResponse(startTime time.Time) *LivenessResponse {
	now := time.Now().UTC()
	uptime := now.Sub(startTime)
	
	return &LivenessResponse{
		Status:        LivenessStatusAlive,
		UptimeSeconds: int64(uptime.Seconds()),
		Timestamp:     now,
	}
}

// IsAlive returns true if the service is alive
func (lr *LivenessResponse) IsAlive() bool {
	return lr.Status == LivenessStatusAlive
}