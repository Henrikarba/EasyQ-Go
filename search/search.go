// Package search provides quantum search functionality using Grover's algorithm.
// It allows developers to search unstructured data without understanding quantum mechanics.
package search

import (
	"errors"

	"github.com/Henrikarba/easyq-go"
	"github.com/Henrikarba/easyq-go/internal/runtime"
)

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

// Options configures the behavior of quantum search operations
type Options struct {
	// MaxAttempts is the maximum number of search attempts to try before giving up.
	// If set to nil, will continue trying until successful or resources are exhausted.
	MaxAttempts *int

	// MaxTargets is the maximum number of target items to sample for the oracle.
	// If set to nil, will use all available targets.
	MaxTargets *int

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
	// This can improve performance when the exact count is known in advance.
	KnownMatchCount *int
}

// Result represents the result of a quantum search operation
type Result struct {
	// Item is the found item
	Item interface{}

	// Index is the position of the item in the original collection
	Index int
}

// DefaultOptions returns a new Options with default values.
func DefaultOptions() Options {
	attempts := 5
	return Options{
		MaxAttempts:           &attempts,
		MaxTargets:            nil,
		IterationStrategy:     Optimal,
		SamplingStrategy:      Auto,
		SampleSize:            100,
		FullScanThreshold:     1000,
		CustomIterationFactor: 1.0,
		CustomIterationOffset: 0,
		EnableLogging:         false,
		KnownMatchCount:       nil,
	}
}

// Search performs a quantum search on the given items using the provided predicate.
// It returns all items that match the predicate and their indices.
// Options may be nil, in which case default options are used.
//
// Example:
//
//	import "github.com/yourusername/easyq-go/search"
//
//	items := []string{"apple", "banana", "cherry", "date"}
//	predicate := func(item string) bool { return len(item) > 5 }
//	results, err := search.Search(items, predicate, nil)
func Search(items interface{}, predicate interface{}, options *Options) ([]Result, error) {
	// Ensure we're initialized
	if err := easyq.EnsureInitialized(); err != nil {
		return nil, err
	}

	// Use default options if none provided
	opts := DefaultOptions()
	if options != nil {
		opts = *options
	}

	// Perform the search
	return runtime.QuantumSearch(items, predicate, opts)
}

// SearchOne performs a quantum search and returns the first matching item.
// This is more efficient than Search when only one result is needed.
func SearchOne(items interface{}, predicate interface{}, options *Options) (*Result, error) {
	// Ensure we're initialized
	if err := easyq.EnsureInitialized(); err != nil {
		return nil, err
	}

	// Use default options if none provided
	opts := DefaultOptions()
	if options != nil {
		opts = *options
	}

	// Set to assume one match exists
	opts.SamplingStrategy = AssumeOne

	// Perform the search
	results, err := runtime.QuantumSearch(items, predicate, opts)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no matching items found")
	}

	return &results[0], nil
}
