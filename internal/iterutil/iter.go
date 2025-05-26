package iterutil

import (
	"maps"
	"reflect"
	"slices"
)

// TODO: use [iter.Seq]

func Permutations(v map[string][]reflect.Value) []map[string]reflect.Value {
	// Initialize the result slice to store all permutations
	var result []map[string]reflect.Value

	keys := slices.Collect(maps.Keys(v))

	// Sort keys for consistent processing order (optional but ensures deterministic output)
	slices.Sort(keys)

	// Define a recursive helper function as a closure
	var generatePermutations func(current map[string]reflect.Value, index int)

	generatePermutations = func(current map[string]reflect.Value, index int) {
		// Base case: if all keys have been processed
		if index == len(keys) {
			// Add the copy to the result
			result = append(result, maps.Clone(current))
			return
		}

		// Get the current key
		key := keys[index]

		// Iterate over all values for the current key
		for _, val := range v[key] {
			// Set the value for the current key
			current[key] = val
			// Recursively generate permutations for the next key
			generatePermutations(current, index+1)
		}
	}

	// Initialize an empty map to build combinations
	current := make(map[string]reflect.Value)

	// Start the recursive generation
	generatePermutations(current, 0)

	return result
}
