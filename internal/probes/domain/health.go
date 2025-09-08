package domain

import "time"

// HealthStatus represents the overall health status
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// CheckStatus represents the status of an individual check
type CheckStatus string

const (
	CheckStatusUp   CheckStatus = "up"
	CheckStatusDown CheckStatus = "down"
)

// Check represents an individual health check result
type Check struct {
	Status         CheckStatus `json:"status"`
	ResponseTimeMs int64       `json:"response_time_ms"`
}

// HealthResponse represents the complete health check response
type HealthResponse struct {
	Status    HealthStatus     `json:"status"`
	Checks    map[string]Check `json:"checks"`
	Timestamp time.Time        `json:"timestamp"`
}

// NewHealthResponse creates a new health response
func NewHealthResponse() *HealthResponse {
	return &HealthResponse{
		Checks:    make(map[string]Check),
		Timestamp: time.Now().UTC(),
	}
}

// AddCheck adds a check result to the health response
func (hr *HealthResponse) AddCheck(name string, status CheckStatus, responseTimeMs int64) {
	hr.Checks[name] = Check{
		Status:         status,
		ResponseTimeMs: responseTimeMs,
	}
}

// DetermineOverallStatus determines the overall health status based on individual checks
func (hr *HealthResponse) DetermineOverallStatus() {
	for _, check := range hr.Checks {
		if check.Status == CheckStatusDown {
			hr.Status = HealthStatusUnhealthy
			return
		}
	}
	hr.Status = HealthStatusHealthy
}

// IsHealthy returns true if the overall status is healthy
func (hr *HealthResponse) IsHealthy() bool {
	return hr.Status == HealthStatusHealthy
}
