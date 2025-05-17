package easyq

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Henrikarba/easyq-go/internal/runtime"
)

// Version information
const (
	// Version is the current version of the EasyQ Go package
	Version = "1.0.0"

	// APIVersion is the version of the API compatible with this package
	APIVersion = "1.0"
)

var (
	// Global initialization state
	initOnce      sync.Once
	initErr       error
	isInitialized bool

	// Global connection configuration
	connectionMutex        sync.RWMutex
	globalConnectionConfig QuantumConnectionConfig
)

// Initialize sets up the EasyQ runtime and prepares it for use.
// It should be called once at the start of your application.
// If not called explicitly, it will be called automatically on first use of
// any functionality in the search or crypto packages.
//
// Returns an error if initialization fails.
func Initialize() error {
	initOnce.Do(func() {
		initErr = runtime.Initialize()
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
		runtime.Shutdown()
		isInitialized = false
	}
}

// GetVersion returns the current version of the EasyQ package.
func GetVersion() string {
	return Version
}

// SetQuantumConnection configures the connection to a quantum computing resource.
// This must be called before using any quantum operations, or the default simulator will be used.
func SetQuantumConnection(config QuantumConnectionConfig) error {
	// Validate the configuration
	if err := validateConnectionConfig(config); err != nil {
		return err
	}

	connectionMutex.Lock()
	defer connectionMutex.Unlock()

	// Store the configuration
	globalConnectionConfig = config

	// If already initialized, we need to reconfigure the connection
	if isInitialized {
		return runtime.ReconfigureConnection(config)
	}

	return nil
}

// GetQuantumConnection returns the current quantum connection configuration
func GetQuantumConnection() QuantumConnectionConfig {
	connectionMutex.RLock()
	defer connectionMutex.RUnlock()
	return globalConnectionConfig
}

// UseDefaultSimulator sets up a simulation backend (no real quantum hardware)
// This is the default if no connection is configured
func UseDefaultSimulator() {
	SetQuantumConnection(QuantumConnectionConfig{
		BackendType: Simulator,
	})
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

// CalculateSecurityMargin calculates the security margin percentage based on the security parameter.
// This is a utility function used by the crypto package.
func CalculateSecurityMargin(securityParameter float64) float64 {
	classicalLimit := 2.0
	quantumMax := 2.0 * math.Sqrt(2.0) // ≈ 2.83

	margin := securityParameter - classicalLimit
	maxMargin := quantumMax - classicalLimit // ≈ 0.83

	if margin <= 0 {
		return 0
	}
	if margin >= maxMargin {
		return 100
	}
	return 100 * margin / maxMargin
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
			return fmt.Errorf("authentication required for %v: provide either Token or Username/Password", config.BackendType)
		}
		return nil
	case LocalQuantumDevice:
		// For local devices, we need an endpoint
		if config.Endpoint == "" {
			return errors.New("endpoint required for LocalQuantumDevice")
		}
		return nil
	case CustomQuantumBackend:
		// For custom backends, we need provider information
		if config.ProviderSettings == nil || config.ProviderSettings["ProviderName"] == "" {
			return errors.New("ProviderName required in ProviderSettings for CustomQuantumBackend")
		}
		return nil
	default:
		return fmt.Errorf("unknown backend type: %v", config.BackendType)
	}
}
