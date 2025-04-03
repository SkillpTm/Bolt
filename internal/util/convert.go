// Package util provides a variation of functions to be used throughout the project
package util

// MakeBoolMap converts a string slice to map[string]true
func MakeBoolMap(input []string) map[string]bool {
	output := make(map[string]bool)

	for _, item := range input {
		output[item] = true
	}

	return output
}
