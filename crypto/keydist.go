package crypto

import (
	"math"

	easyq "github.com/Henrikarba/easyq-go"
	"github.com/Henrikarba/easyq-go/bridge"
)

// DefaultKeyDistributionOptions returns a new set of default options for key distribution
func DefaultKeyDistributionOptions() easyq.KeyDistributionOptions {
	return easyq.KeyDistributionOptions{
		KeyLength:              256,
		SecurityLevel:          3,
		SecurityThreshold:      2.2,
		MaxAttempts:            5,
		EnableLogging:          false,
		AuthenticationMode:     easyq.Standard,
		EnableErrorCorrection:  true,
		MaxAcceptableErrorRate: 0.12,
	}
}

// GenerateKey generates a cryptographically secure key using quantum key distribution.
// This uses the E91 protocol with entangled quantum particles to create a secure key
// that is protected by the fundamental laws of quantum physics.
//
// Example:
//
//	// Generate a 256-bit quantum-secure key with default options
//	result, err := crypto.GenerateKey(nil)
func GenerateKey(options *easyq.KeyDistributionOptions) (*easyq.KeyDistributionResult, error) {
	// Ensure we're initialized
	if err := easyq.EnsureInitialized(); err != nil {
		return nil, err
	}

	// Use default options if none provided
	opts := DefaultKeyDistributionOptions()
	if options != nil {
		opts = *options
	}

	// Validate options
	if opts.KeyLength <= 0 {
		return nil, easyq.ErrInvalidLength
	}

	if opts.SecurityLevel < 1 || opts.SecurityLevel > 5 {
		return nil, easyq.ErrInvalidSecurityLevel
	}

	// Generate key using the bridge
	rawResult, err := bridge.GenerateKey(opts)
	if err != nil {
		return nil, err
	}

	// Convert raw result to KeyDistributionResult
	result := &easyq.KeyDistributionResult{}

	// Process success flag
	if success, ok := rawResult["Success"].(bool); ok {
		result.Success = success
	}

	// Process key data if successful
	if result.Success {
		// Process key bytes
		if keyData, ok := rawResult["Key"].([]interface{}); ok {
			result.Key = make([]byte, len(keyData))
			for i, b := range keyData {
				if bFloat, ok := b.(float64); ok {
					result.Key[i] = byte(bFloat)
				}
			}
		}

		// Process authentication tag
		if tagData, ok := rawResult["AuthenticationTag"].([]interface{}); ok {
			result.AuthenticationTag = make([]byte, len(tagData))
			for i, b := range tagData {
				if bFloat, ok := b.(float64); ok {
					result.AuthenticationTag[i] = byte(bFloat)
				}
			}
		}
	}

	// Process security parameter
	if secParam, ok := rawResult["SecurityParameter"].(float64); ok {
		result.SecurityParameter = secParam
	}

	// Process error rate
	if errorRate, ok := rawResult["ErrorRate"].(float64); ok {
		result.ErrorRate = errorRate
	}

	// Process entangled pairs
	if pairs, ok := rawResult["EntangledPairsCreated"].(float64); ok {
		result.EntangledPairsCreated = int(pairs)
	}

	// Process failure reason if not successful
	if !result.Success {
		if reason, ok := rawResult["FailureReason"].(string); ok {
			result.FailureReason = reason
		} else {
			result.FailureReason = "Unknown failure"
		}
		return result, easyq.ErrKeyGenerationFailed
	}

	return result, nil
}

// VerifyChannelSecurity checks if a quantum channel is secure for key distribution
// without actually generating a full key. This can be used to detect eavesdropping.
//
// Returns:
// - isSecure: whether the channel appears secure (no eavesdropping detected)
// - securityParameter: the calculated security parameter (CHSH value)
// - errorRate: the observed error rate in measurements
// - error: any error that occurred during verification
func VerifyChannelSecurity(options *easyq.KeyDistributionOptions) (bool, float64, float64, error) {
	// Use a small key length for verification only
	opts := DefaultKeyDistributionOptions()
	if options != nil {
		opts = *options
	}

	// Override certain options for verification
	opts.KeyLength = 32  // Smaller key for verification only
	opts.MaxAttempts = 2 // Fewer attempts since we're just testing

	// Generate a short key to test the channel
	result, err := GenerateKey(&opts)
	if err != nil {
		// If it's not specifically a key generation failure, return the error
		if err != easyq.ErrKeyGenerationFailed {
			return false, 0, 0, err
		}
		// Otherwise, we'll continue and return the security parameters even though generation failed
	}

	isSecure := result.Success &&
		result.SecurityParameter >= opts.SecurityThreshold &&
		result.ErrorRate <= opts.MaxAcceptableErrorRate

	return isSecure, result.SecurityParameter, result.ErrorRate, nil
}

// CalculateSecurityMargin calculates the security margin percentage based on the security parameter.
// The security parameter (CHSH value) ranges from 2.0 (classical limit) to 2.83 (quantum maximum).
// The margin is expressed as a percentage of how far above the classical limit the parameter is.
//
// Example:
//
//	// Calculate the security margin for a CHSH value of 2.5
//	margin := crypto.CalculateSecurityMargin(2.5)
//	// Returns approximately 60.2% (security margin)
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
