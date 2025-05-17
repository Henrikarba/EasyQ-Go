// Package runtime manages the EasyQ runtime and provides access to quantum operations.
// This is an internal package not meant to be used directly by users.
package runtime

// #cgo LDFLAGS: -L${SRCDIR}/../../lib/${GOOS}_${GOARCH} -leasyq_bridge
// #include <stdlib.h>
// #include "../../internal/bridge/bridge.h"
import "C"
import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

var (
	// Global initialization state
	runtimeMutex  sync.Mutex
	isInitialized bool
)

// Initialize sets up the EasyQ runtime and prepares it for use.
func Initialize() error {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if isInitialized {
		return nil
	}

	result := C.EasyQ_Initialize()
	if result != 0 {
		return errors.New("failed to initialize EasyQ runtime")
	}

	// Set up a finalizer to ensure Shutdown is called when the program exits
	runtime.SetFinalizer(&finalizerObject, func(_ *int) {
		Shutdown()
	})

	isInitialized = true
	return nil
}

// Shutdown cleans up resources used by the EasyQ runtime.
func Shutdown() {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if isInitialized {
		C.EasyQ_Shutdown()
		isInitialized = false
	}
}

// ReconfigureConnection reconfigures the connection to a quantum computing resource.
func ReconfigureConnection(config interface{}) error {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if !isInitialized {
		return errors.New("runtime not initialized")
	}

	// For now, this is a placeholder
	// In the actual implementation, we would serialize the config and pass it to the C layer
	return nil
}

// QuantumSearch performs a quantum search operation.
func QuantumSearch(items interface{}, predicate interface{}, options interface{}) ([]interface{}, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, we would:
	// 1. Serialize the items, predicate, and options
	// 2. Pass them to the C# layer via CGo
	// 3. Deserialize the results
	return nil, errors.New("not implemented yet")
}

// QuantumRNG generates a random integer using quantum measurement.
func QuantumRNG(min, max int) (int, error) {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if !isInitialized {
		return 0, errors.New("runtime not initialized")
	}

	var result C.int
	status := C.EasyQ_GenerateRandomInt(C.int(min), C.int(max), &result)
	if status != 0 {
		return 0, errors.New("failed to generate random integer")
	}

	return int(result), nil
}

// QuantumRandomBytes generates random bytes using quantum measurement.
func QuantumRandomBytes(length int) ([]byte, error) {
	runtimeMutex.Lock()
	defer runtimeMutex.Unlock()

	if !isInitialized {
		return nil, errors.New("runtime not initialized")
	}

	// Allocate a buffer for the result
	buffer := C.malloc(C.size_t(length))
	if buffer == nil {
		return nil, errors.New("failed to allocate memory")
	}
	defer C.free(buffer)

	status := C.EasyQ_GenerateRandomBytes(C.int(length), (*C.uchar)(buffer))
	if status != 0 {
		return nil, errors.New("failed to generate random bytes")
	}

	// Copy the result to a Go slice
	result := make([]byte, length)
	C.memcpy(unsafe.Pointer(&result[0]), buffer, C.size_t(length))

	return result, nil
}

// Dummy object used for finalizer
var finalizerObject = new(int)
