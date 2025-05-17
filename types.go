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
