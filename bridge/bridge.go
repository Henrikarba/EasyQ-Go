// Package bridge provides direct communication with the native quantum operations
// via a DLL/shared library interface.
package bridge

// #cgo windows LDFLAGS: -L${SRCDIR}/../../lib/windows_amd64 -lEasyQBridge
// #cgo linux LDFLAGS: -L${SRCDIR}/../../lib/linux_amd64 -lEasyQBridge
// #cgo darwin LDFLAGS: -L${SRCDIR}/../../lib/darwin_amd64 -lEasyQBridge
// #include <stdlib.h>
// #include <stdint.h>
// #include "bridge.h"

import (
	"C"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

// Status codes from the DLL
const (
	StatusSuccess              = 0
	StatusErrorGeneral         = 1
	StatusErrorNotInitialized  = 2
	StatusErrorInvalidArgument = 3
	StatusErrorRuntime         = 4
	StatusErrorTimeout         = 5
)

var (
	bridgeMutex   sync.Mutex
	isInitialized bool
)

// Initialize initializes the quantum bridge.
func Initialize() error {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if isInitialized {
		return nil
	}

	result := C.EasyQ_Initialize()
	if result != StatusSuccess {
		return fmt.Errorf("failed to initialize quantum bridge: error code %d", result)
	}

	isInitialized = true
	return nil
}

// Shutdown cleans up resources used by the quantum bridge.
func Shutdown() {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if isInitialized {
		C.EasyQ_Shutdown()
		isInitialized = false
	}
}

// ConfigureConnection configures the connection to a quantum computing resource.
func ConfigureConnection(config interface{}) error {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if !isInitialized {
		return errors.New("bridge not initialized")
	}

	// Convert the config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal connection config: %w", err)
	}

	// Convert JSON to C string
	cConfigJSON := C.CString(string(configJSON))
	defer C.free(unsafe.Pointer(cConfigJSON))

	// Call the DLL function
	result := C.EasyQ_ConfigureConnection(cConfigJSON)
	if result != StatusSuccess {
		return fmt.Errorf("failed to configure quantum connection: error code %d", result)
	}

	return nil
}

// Search performs a quantum search using Grover's algorithm.
func Search(items interface{}, predicate interface{}, options interface{}) ([]interface{}, error) {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if !isInitialized {
		return nil, errors.New("bridge not initialized")
	}

	// Convert parameters to JSON
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}

	predicateJSON, err := json.Marshal(predicate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal predicate: %w", err)
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}

	// Convert JSON to C strings
	cItemsJSON := C.CString(string(itemsJSON))
	defer C.free(unsafe.Pointer(cItemsJSON))

	cPredicateJSON := C.CString(string(predicateJSON))
	defer C.free(unsafe.Pointer(cPredicateJSON))

	cOptionsJSON := C.CString(string(optionsJSON))
	defer C.free(unsafe.Pointer(cOptionsJSON))

	// Prepare for result
	var cResultJSON *C.char

	// Call the DLL function
	result := C.EasyQ_Search(cItemsJSON, cPredicateJSON, cOptionsJSON, &cResultJSON)
	if result != StatusSuccess {
		return nil, fmt.Errorf("quantum search failed: error code %d", result)
	}

	// Convert result back to Go and free the C string
	goResultJSON := C.GoString(cResultJSON)
	C.EasyQ_FreeString(cResultJSON)

	// Unmarshal the result
	var searchResults []interface{}
	err = json.Unmarshal([]byte(goResultJSON), &searchResults)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return searchResults, nil
}

// GenerateRandomInt generates a random integer using quantum measurement.
func GenerateRandomInt(min, max int) (int, error) {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if !isInitialized {
		return 0, errors.New("bridge not initialized")
	}

	// Prepare for result
	var result C.int

	// Call the DLL function
	status := C.EasyQ_GenerateRandomInt(C.int(min), C.int(max), &result)
	if status != StatusSuccess {
		return 0, fmt.Errorf("quantum RNG failed: error code %d", status)
	}

	return int(result), nil
}

// GenerateRandomBytes generates random bytes using quantum measurement.
func GenerateRandomBytes(length int) ([]byte, error) {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if !isInitialized {
		return nil, errors.New("bridge not initialized")
	}

	// Allocate a buffer for the result
	buffer := make([]byte, length)

	// Call the DLL function
	status := C.EasyQ_GenerateRandomBytes(C.int(length), (*C.uchar)(unsafe.Pointer(&buffer[0])))
	if status != StatusSuccess {
		return nil, fmt.Errorf("quantum random bytes generation failed: error code %d", status)
	}

	return buffer, nil
}

// GenerateKey generates a key using quantum key distribution.
func GenerateKey(options interface{}) (map[string]interface{}, error) {
	bridgeMutex.Lock()
	defer bridgeMutex.Unlock()

	if !isInitialized {
		return nil, errors.New("bridge not initialized")
	}

	// Convert options to JSON
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}

	// Convert JSON to C string
	cOptionsJSON := C.CString(string(optionsJSON))
	defer C.free(unsafe.Pointer(cOptionsJSON))

	// Prepare for result
	var cResultJSON *C.char

	// Call the DLL function
	status := C.EasyQ_GenerateKey(cOptionsJSON, &cResultJSON)
	if status != StatusSuccess {
		return nil, fmt.Errorf("quantum key distribution failed: error code %d", status)
	}

	// Convert result back to Go and free the C string
	goResultJSON := C.GoString(cResultJSON)
	C.EasyQ_FreeString(cResultJSON)

	// Unmarshal the result
	var keyResult map[string]interface{}
	err = json.Unmarshal([]byte(goResultJSON), &keyResult)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal key distribution result: %w", err)
	}

	return keyResult, nil
}
