package easyq

// QuantumBackendType defines the type of quantum backend to use
type QuantumBackendType int

const (
	// Simulator runs quantum operations on a classical simulator
	Simulator QuantumBackendType = iota

	// MicrosoftQuantumCloud connects to Microsoft's quantum cloud service
	MicrosoftQuantumCloud

	// IBMQuantumExperience connects to IBM's quantum experience
	IBMQuantumExperience

	// GoogleQuantumAI connects to Google's quantum AI platform
	GoogleQuantumAI

	// LocalQuantumDevice connects to a local quantum computing device
	LocalQuantumDevice

	// CustomQuantumBackend connects to a custom quantum backend
	CustomQuantumBackend
)

// QuantumConnectionConfig holds information needed to connect to a quantum computing resource
type QuantumConnectionConfig struct {
	// The type of quantum backend to use
	BackendType QuantumBackendType

	// Connection information for remote quantum computers
	Endpoint string // IP address or hostname of the quantum computer
	Port     int    // Port number for connection
	Region   string // Region for cloud services
	Username string // Authentication username
	Password string // Authentication password or API key
	Token    string // Alternative authentication token

	// Custom provider-specific settings
	ProviderSettings map[string]string
}

// SearchResult represents the result of a quantum search operation
type SearchResult struct {
	// Item is the found item
	Item interface{}

	// Index is the position of the item in the original collection
	Index int
}

// SearchOptions configures the behavior of quantum search operations
type SearchOptions struct {
	// MaxAttempts is the maximum number of search attempts to try before giving up.
	// If set to 0, will use the default (5).
	MaxAttempts int

	// MaxTargets is the maximum number of target items to sample for the oracle.
	// If set to 0, will use all available targets.
	MaxTargets int

	// IterationStrategy determines how to calculate the number of Grover iterations.
	IterationStrategy IterationStrategy

	// SamplingStrategy determines how to estimate the number of matches.
	SamplingStrategy SamplingStrategy

	// SampleSize is the number of samples to use when estimating match counts.
	// Only used with SamplingStrategy.Sampling.
	SampleSize int

	// FullScanThreshold is the maximum database size to scan completely.
	// Only used with SamplingStrategy.Auto.
	FullScanThreshold int

	// CustomIterationFactor multiplies the standard optimal iteration count.
	// Only used when IterationStrategy is set to Custom.
	CustomIterationFactor float64

	// CustomIterationOffset adds to the iteration count.
	// Only used when IterationStrategy is set to Custom.
	CustomIterationOffset int

	// EnableLogging determines whether to log detailed information about the quantum search process.
	EnableLogging bool

	// KnownMatchCount is an optional parameter to specify the exact number of matching items.
	// If set to 0, will be ignored.
	KnownMatchCount int
}

// IterationStrategy defines strategies for determining the number of Grover iterations
type IterationStrategy int

const (
	// Optimal is standard optimal iterations (PI / (4 * angle) - 0.5)
	Optimal IterationStrategy = iota

	// SingleIteration is one iteration only (useful for small search spaces)
	SingleIteration

	// Aggressive is a more aggressive approach (PI / (4 * angle))
	Aggressive

	// Conservative is a conservative approach (PI / (4 * angle) - 1)
	Conservative

	// HalfOptimal is half the standard iterations (PI / (8 * angle))
	HalfOptimal

	// Custom is a custom iteration count using CustomIterationFactor and CustomIterationOffset
	Custom
)

// SamplingStrategy defines strategies for estimating the number of matches in the database
type SamplingStrategy int

const (
	// Auto automatically chooses between FullScan and Sampling based on database size
	Auto SamplingStrategy = iota

	// FullScan always scans the entire database (accurate but may be slow for large databases)
	FullScan

	// Sampling uses random sampling to estimate (faster but less accurate)
	Sampling

	// AssumeOne assumes only one match exists (fastest but only appropriate when you know there's exactly one match)
	AssumeOne

	// UserProvided uses a specific count provided by the user
	UserProvided
)

// KeyDistributionOptions configures the behavior of quantum key distribution operations
type KeyDistributionOptions struct {
	// KeyLength is the desired length of the generated key in bits.
	KeyLength int

	// SecurityLevel is the security level (1-5). Higher values increase security but decrease efficiency.
	SecurityLevel int

	// SecurityThreshold is the minimum security parameter (CHSH value) required to consider a channel secure.
	// Classical limit is 2.0. Quantum maximum is 2√2 ≈ 2.83.
	// Recommended values: 2.2-2.4 for standard security.
	SecurityThreshold float64

	// MaxAttempts is maximum number of attempts to generate a key before giving up.
	// If set to 0, will use the default (5).
	MaxAttempts int

	// EnableLogging determines whether to log detailed information about the key distribution process.
	EnableLogging bool

	// AuthenticationMode is the authentication mode to use for key verification.
	AuthenticationMode AuthenticationMode

	// PreSharedSecret is an optional custom pre-shared authentication secret.
	// If nil, a random one will be generated.
	PreSharedSecret []byte

	// EnableErrorCorrection determines if error correction should be performed on the raw key.
	EnableErrorCorrection bool

	// MaxAcceptableErrorRate is the maximum acceptable error rate before aborting key generation.
	MaxAcceptableErrorRate float64
}

// AuthenticationMode defines the authentication mode for quantum key distribution
type AuthenticationMode int

const (
	// None means no authentication is performed
	None AuthenticationMode = iota

	// Standard authentication uses pre-shared secret
	Standard

	// Enhanced authentication with additional quantum verification
	Enhanced
)

// KeyDistributionResult represents the result of a quantum key distributi
