package easyq

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrNotInitialized is returned when operations are performed without initialization
	ErrNotInitialized = errors.New("easyq: not initialized")

	// ErrInvalidAuth is returned when authentication details are missing
	ErrInvalidAuth = errors.New("easyq: missing authentication details for selected backend")

	// ErrMissingEndpoint is returned when endpoint is missing for local device
	ErrMissingEndpoint = errors.New("easyq: missing endpoint for local device")

	// ErrMissingProviderInfo is returned when provider information is missing
	ErrMissingProviderInfo = errors.New("easyq: missing provider information for custom backend")

	// ErrUnknownBackend is returned when an unknown backend type is specified
	ErrUnknownBackend = errors.New("easyq: unknown backend type")

	// ErrNoMatches is returned when a search operation finds no matches
	ErrNoMatches = errors.New("easyq: no matching items found")

	// ErrInvalidRange is returned when an invalid range is specified
	ErrInvalidRange = errors.New("easyq: invalid range (min must be less than max)")

	// ErrInvalidLength is returned when an invalid length is specified
	ErrInvalidLength = errors.New("easyq: invalid length (must be greater than zero)")

	// ErrInvalidSecurityLevel is returned when an invalid security level is specified
	ErrInvalidSecurityLevel = errors.New("easyq: invalid security level (must be 1-5)")

	// ErrKeyGenerationFailed is returned when key generation fails
	ErrKeyGenerationFailed = errors.New("easyq: key generation failed")
)

// BridgeError represents an error from the native bridge
type BridgeError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *BridgeError) Error() string {
	return fmt.Sprintf("easyq bridge error (code %d): %s", e.Code, e.Message)
}

// NewBridgeError creates a new BridgeError
func NewBridgeError(code int, message string) *BridgeError {
	return &BridgeError{
		Code:    code,
		Message: message,
	}
}
