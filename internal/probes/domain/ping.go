package domain

// PingResponse represents the response for a ping request
type PingResponse struct {
	Message string `json:"message"`
}

// NewPingResponse creates a new ping response with the standard "PONG" message
func NewPingResponse() *PingResponse {
	return &PingResponse{
		Message: "PONG",
	}
}