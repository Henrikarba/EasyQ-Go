// Package easyq provides a developer-friendly API for quantum computing operations
// without requiring specialized knowledge of quantum mechanics or computing principles.
package easyq

import (
	"sync"

	"github.com/Henrikarba/easyq-go/bridge"
)

var (
	// Global initialization state
	initOnce      sync.Once
	initErr       error
	isInitialized bool
)

// Initialize sets up the EasyQ runtime and prepares it for use.
// It should be called once at the start of your application.
// If not called explicitly, it will be called automatically on first use of
// any functionality in the search or crypto packages.
//
// Returns an error if initialization fails.
func Initialize() error {
	initOnce.Do(func() {
		initErr = bridge.Initialize()
		if initErr == nil {
			isInitialized = true
		}
	})
	return initErr
}

// Shutdown cleans up resources used by the EasyQ runtime.
// It should be called when your application is shutting down.
func Shutdown() {
	if isInitialized {
		bridge.Shutdown()
		isInitialized = false
	}
}

// IsInitialized returns whether the EasyQ runtime has been initialized
func IsInitialized() bool {
	return isInitialized
}

// EnsureInitialized ensures the package is initialized
// This is used internally by the subpackages
func EnsureInitialized() error {
	if !isInitialized {
		return Initialize()
	}
	return nil
}

// UseDefaultSimulator sets up a simulation backend (no real quantum hardware)
// This is the default if no connection is configured
func UseDefaultSimulator() error {
	return SetQuantumConnection(QuantumConnectionConfig{
		BackendType: Simulator,
	})
}

// SetQuantumConnection configures the connection to a quantum computing resource.
// This must be called before using any quantum operations, or the default simulator will be used.
func SetQuantumConnection(config QuantumConnectionConfig) error {
	// Validate the configuration
	if err := validateConnectionConfig(config); err != nil {
		return err
	}

	// Ensure we're initialized
	if err := EnsureInitialized(); err != nil {
		return err
	}

	// Configure the connection through the bridge
	return bridge.ConfigureConnection(config)
}

// GetVersion returns the current version of the EasyQ package.
func GetVersion() string {
	return "0.1.0"
}

// Validate connection configuration
func validateConnectionConfig(config QuantumConnectionConfig) error {
	// Check for invalid combinations
	switch config.BackendType {
	case Simulator:
		// No additional validation needed for simulator
		return nil
	case MicrosoftQuantumCloud, IBMQuantumExperience, GoogleQuantumAI:
		// For cloud services, we need authentication
		if config.Token == "" && (config.Username == "" || config.Password == "") {
			return ErrInvalidAuth
		}
		return nil
	case LocalQuantumDevice:
		// For local devices, we need an endpoint
		if config.Endpoint == "" {
			return ErrMissingEndpoint
		}
		return nil
	case CustomQuantumBackend:
		// For custom backends, we need provider information
		if config.ProviderSettings == nil || config.ProviderSettings["ProviderName"] == "" {
			return ErrMissingProviderInfo
		}
		return nil
	default:
		return ErrUnknownBackend
	}
}
