package dto

// AuthRequest represents the authentication request payload
type AuthRequest struct {
	Secret string `json:"secret" validate:"required" example:"your-secret-key"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt int64  `json:"expires_at" example:"1642492800"`
}
