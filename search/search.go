// Package search provides quantum search functionality using Grover's algorithm.
// It allows developers to search unstructured data without understanding quantum mechanics.
package search

import (
	"errors"
	"fmt"
	"reflect"

	easyq "github.com/Henrikarba/easyq-go"
	"github.com/Henrikarba/easyq-go/bridge"
)

// DefaultOptions returns a new Options with default values.
func DefaultOptions() easyq.SearchOptions {
	return easyq.SearchOptions{
		MaxAttempts:           5,
		IterationStrategy:     easyq.Optimal,
		SamplingStrategy:      easyq.Auto,
		SampleSize:            100,
		FullScanThreshold:     1000,
		CustomIterationFactor: 1.0,
		CustomIterationOffset: 0,
		EnableLogging:         false,
	}
}

// Search performs a quantum search on the given items using the provided predicate.
// It returns all items that match the predicate and their indices.
// Options may be nil, in which case default options are used.
//
// Example:
//
//	items := []string{"apple", "banana", "cherry", "date"}
//	predicate := func(item string) bool { return len(item) > 5 }
//	results, err := search.Search(items, predicate, nil)
func Search(items interface{}, predicate interface{}, options *easyq.SearchOptions) ([]easyq.SearchResult, error) {
	// Validate inputs
	if err := validateInputs(items, predicate); err != nil {
		return nil, err
	}

	// Ensure we're initialized
	if err := easyq.EnsureInitialized(); err != nil {
		return nil, err
	}

	// Use default options if none provided
	opts := DefaultOptions()
	if options != nil {
		opts = *options
	}

	// Prepare mapped predicate for serialization
	mappedPredicate, err := convertPredicate(predicate, reflect.TypeOf(items).Elem())
	if err != nil {
		return nil, err
	}

	// Perform the search through the bridge
	rawResults, err := bridge.Search(items, mappedPredicate, opts)
	if err != nil {
		return nil, err
	}

	// Convert raw results to SearchResult objects
	results := make([]easyq.SearchResult, 0, len(rawResults))
	for _, rawResult := range rawResults {
		resultMap, ok := rawResult.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected result format: %T", rawResult)
		}

		var result easyq.SearchResult

		// Extract index
		indexValue, ok := resultMap["Index"]
		if !ok {
			return nil, errors.New("result missing Index field")
		}
		index, ok := indexValue.(float64)
		if !ok {
			return nil, fmt.Errorf("unexpected index type: %T", indexValue)
		}
		result.Index = int(index)

		// Extract item
		itemValue, ok := resultMap["Item"]
		if !ok {
			return nil, errors.New("result missing Item field")
		}
		result.Item = itemValue

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, easyq.ErrNoMatches
	}

	return results, nil
}

// SearchOne performs a quantum search and returns the first matching item.
// This is more efficient than Search when only one result is needed.
func SearchOne(items interface{}, predicate interface{}, options *easyq.SearchOptions) (*easyq.SearchResult, error) {
	// Use default options if none provided
	opts := DefaultOptions()
	if options != nil {
		opts = *options
	}

	// Set to assume one match exists
	opts.SamplingStrategy = easyq.AssumeOne
	opts.MaxAttempts = 3 // Less attempts since we only need one match

	// Perform the search
	results, err := Search(items, predicate, &opts)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, easyq.ErrNoMatches
	}

	return &results[0], nil
}

// validateInputs checks that the items and predicate are valid for quantum search
func validateInputs(items interface{}, predicate interface{}) error {
	// Check if items is a slice or array
	itemsType := reflect.TypeOf(items)
	if itemsType == nil || (itemsType.Kind() != reflect.Slice && itemsType.Kind() != reflect.Array) {
		return errors.New("items must be a slice or array")
	}

	// Check if predicate is a function
	predicateType := reflect.TypeOf(predicate)
	if predicateType == nil || predicateType.Kind() != reflect.Func {
		return errors.New("predicate must be a function")
	}

	// Check if predicate has correct signature (func(T) bool)
	if predicateType.NumIn() != 1 || predicateType.NumOut() != 1 {
		return errors.New("predicate must have signature func(T) bool")
	}

	// Check if predicate input type matches items element type
	if !predicateType.In(0).AssignableTo(itemsType.Elem()) {
		return errors.New("predicate input type must match items element type")
	}

	// Check if predicate output type is bool
	if predicateType.Out(0).Kind() != reflect.Bool {
		return errors.New("predicate must return a boolean")
	}

	return nil
}

// convertPredicate converts a Go function to a format that can be serialized
// and understood by the bridge implementation
func convertPredicate(predicate interface{}, itemType reflect.Type) (interface{}, error) {
	// For simplicity, we'll create a placeholder representation here
	// In a real implementation, you might serialize the function logic or use a different approach

	predicateValue := reflect.ValueOf(predicate)
	predicateType := predicateValue.Type()

	// Create a representation of the predicate
	return map[string]interface{}{
		"Type":           "Function",
		"InputType":      itemType.String(),
		"ReturnType":     predicateType.Out(0).String(),
		"SerializedFunc": fmt.Sprintf("%v", predicate),
	}, nil
}
