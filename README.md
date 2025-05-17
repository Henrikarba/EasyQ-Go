# EasyQ for Go

A developer-friendly quantum computing framework for Go, allowing regular developers to leverage quantum algorithms without specialized knowledge of quantum mechanics.

## Overview

EasyQ provides a clean, idiomatic Go interface to quantum computing operations, abstracting away the complexity of quantum operations. It supports:

- **Quantum Search**: Implementation of Grover's algorithm for searching unstructured data
- **Quantum Random Number Generation**: True random number generation using quantum properties
- **Quantum Key Distribution**: Secure communications using quantum principles

The package works with both quantum simulators and real quantum hardware (when available) through a unified API.

## Installation

```bash
go get github.com/yourusername/easyq-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/easyq-go"
    "github.com/yourusername/easyq-go/search"
    "github.com/yourusername/easyq-go/crypto"
)

func main() {
    // Initialize with the default simulator
    easyq.UseDefaultSimulator()

    // Example 1: Quantum Search
    people := []string{"Alice", "Bob", "Charlie", "David", "Emma", "Frank"}

    // Define a predicate to search for people whose names start with "E"
    results, err := search.Search(people, func(name string) bool {
        return len(name) > 0 && name[0] == 'E'
    }, nil)

    if err != nil {
        log.Fatalf("Search error: %v", err)
    }

    for _, result := range results {
        fmt.Printf("Found: %s at index %d\n", result.Item, result.Index)
    }

    // Example 2: Quantum Random Number Generation
    number, err := crypto.RandomInt(1, 100)
    if err != nil {
        log.Fatalf("Random number error: %v", err)
    }
    fmt.Printf("Quantum random number: %d\n", number)

    // Clean up resources when done
    easyq.Shutdown()
}
```

## Connecting to Quantum Hardware

```go
package main

import (
    "fmt"
    "log"

    "github.com/yourusername/easyq-go"
    "github.com/yourusername/easyq-go/crypto"
)

func main() {
    // Configure connection to a quantum computer
    err := easyq.SetQuantumConnection(easyq.QuantumConnectionConfig{
        BackendType: easyq.IBMQuantumExperience,
        Username:    "your-username",
        Token:       "your-api-token",
        ProviderSettings: map[string]string{
            "Device": "ibmq_manila",
        },
    })

    if err != nil {
        log.Fatalf("Connection error: %v", err)
    }

    // Generate true quantum random bytes
    randomBytes, err := crypto.RandomBytes(32)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Printf("Generated %d quantum-secure random bytes\n", len(randomBytes))
}
```

## Documentation

For detailed documentation, see the [Go package documentation](https://pkg.go.dev/github.com/yourusername/easyq-go).

### Key Features

- **Pure Quantum Logic**: All operations use true quantum algorithms, without classical fallbacks
- **Developer-Friendly API**: No knowledge of quantum mechanics required
- **Simulation & Hardware Support**: Works with simulators or real quantum hardware
- **Production-Ready**: Built for reliability and performance in real-world applications

### Package Structure

- **github.com/yourusername/easyq-go**: Main package with initialization and configuration
- **github.com/yourusername/easyq-go/search**: Quantum search functionality
- **github.com/yourusername/easyq-go/crypto**: Quantum cryptography functions

## License

[MIT License](LICENSE)
