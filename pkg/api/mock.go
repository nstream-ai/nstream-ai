package api

import (
	"fmt"
	"time"
)

// MockAuthResponse represents a mock response from the auth service
type MockAuthResponse struct {
	Valid     bool      `json:"valid"`
	ExpiresAt time.Time `json:"expires_at"`
	Error     string    `json:"error,omitempty"`
}

// MockValidateToken simulates a gRPC call to validate an auth token
func MockValidateToken(token string) (*MockAuthResponse, error) {
	// Simulate network delay
	time.Sleep(500 * time.Millisecond)

	// Check if token is empty
	if token == "" {
		return &MockAuthResponse{
			Valid: false,
			Error: "empty token",
		}, nil
	}

	// Check if token is expired (mock tokens expire after 24 hours)
	// For testing, we'll consider tokens starting with "expired-" as expired
	if token[:8] == "expired-" {
		return &MockAuthResponse{
			Valid:     false,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Error:     "token expired",
		}, nil
	}

	// Valid token
	return &MockAuthResponse{
		Valid:     true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// MockValidateClusterToken simulates a gRPC call to validate a cluster token
func MockValidateClusterToken(token string) (*MockAuthResponse, error) {
	// Simulate network delay
	time.Sleep(500 * time.Millisecond)

	// Check if token is empty
	if token == "" {
		return &MockAuthResponse{
			Valid: false,
			Error: "empty token",
		}, nil
	}

	// Check if token is expired (mock tokens expire after 24 hours)
	// For testing, we'll consider tokens starting with "expired-" as expired
	if token[:8] == "expired-" {
		return &MockAuthResponse{
			Valid:     false,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Error:     "token expired",
		}, nil
	}

	// Valid token
	return &MockAuthResponse{
		Valid:     true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// MockValidateUser simulates a gRPC call to validate a user
func MockValidateUser(email string) (bool, error) {
	// Simulate network delay
	time.Sleep(500 * time.Millisecond)

	// Check if email is empty
	if email == "" {
		return false, fmt.Errorf("empty email")
	}

	// For testing, consider any non-empty email as valid
	return true, nil
}
