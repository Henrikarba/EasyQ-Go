# EasyQ for Go

A developer-friendly quantum computing framework for Go, allowing regular developers to leverage quantum algorithms without specialized knowledge of quantum mechanics.

## Overview

EasyQ provides a clean, idiomatic Go interface to quantum computing operations, abstracting away the complexity of quantum operations. It supports:

- **Quantum Search**: Implementation of Grover's algorithm for searching unstructured data
- **Quantum Random Number Generation**: True random number generation using quantum properties
- **Quantum Key Distribution**: Secure communications using quantum principles

The package works with both quantum simulators and real quantum hardware (when available) through a unified API.

## Architecture

EasyQ for Go uses a layered architecture:

1. **Go API Layer**: Clean, idiomatic Go interface (`github.com/yourusername/easyq-go`)
2. **Native Bridge**: Native shared library (.dll/.so/.dylib) built from the C# bridge
3. **C#/Q# Backend**: Core quantum operations implemented in C# and Q#

This approach gives you the best of both worlds - easy integration for Go developers and the power of Microsoft's Quantum Development Kit.

## Installation

```bash
go get github.com/yourusername/easyq-go
```

### Prerequisites

EasyQ for Go requires the EasyQBridge shared library which is included in the distribution:

- Windows: `lib/windows_amd64/EasyQBridge.dll`
- Linux: `lib/linux_amd64/libEasyQBridge.so`
- macOS: `lib/darwin_amd64/libEasyQBridge.dylib`

The appropriate library for your platform is automatically loaded at runtime.

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

## Building from Source

### Building the Go Package

```bash
git clone https://github.com/yourusername/easyq-go.git
cd easyq-go
go build ./...
```

### Building the Native Bridge

The native bridge needs to be built from the C# source code:

```bash
cd src
dotnet publish Bridge/Bridge.csproj -c Release -r win-x64
dotnet publish Bridge/Bridge.csproj -c Release -r linux-x64
dotnet publish Bridge/Bridge.csproj -c Release -r osx-x64
```

Copy the generated binaries to the appropriate `lib` directories:

```
lib/
├── windows_amd64/
│   └── EasyQBridge.dll
├── linux_amd64/
│   └── libEasyQBridge.so
└── darwin_amd64/
    └── libEasyQBridge.dylib
```

## Documentation

### Key Features

- **Pure Quantum Logic**: All operations use true quantum algorithms, without classical fallbacks
- **Developer-Friendly API**: No knowledge of quantum mechanics required
- **Simulation & Hardware Support**: Works with simulators or real quantum hardware
- **Production-Ready**: Built for reliability and performance in real-world applications

### Package Structure

- **github.com/yourusername/easyq-go**: Main package with initialization and configuration
- **github.com/yourusername/easyq-go/search**: Quantum search functionality
- **github.com/yourusername/easyq-go/crypto**: Quantum cryptography functions

## How It Works

When you use the EasyQ Go package:

1. The Go code calls into the native bridge using CGo
2. The native bridge (C# compiled to native code) processes the request
3. The C# code calls the Q# quantum operations
4. Results are passed back through the bridge to Go

This multi-layer approach allows for:

- High performance (native code)
- Clean integration with Go
- Full quantum capabilities via Q#

## License

[MIT License](LICENSE)
