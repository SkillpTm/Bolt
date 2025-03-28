// Package sslslices includes functions that make the use of slices easier for me
package sslslices

// <---------------------------------------------------------------------------------------------------->

// Contains checks if the element can be found inside of the slice
func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}

	return false
}

// ConvertInterface converts a []interface{} to a slice of any type
func ConvertInterface[T any](sliceInput []interface{}) []T {
	newSlice := []T{}

	for _, item := range sliceInput {
		newSlice = append(newSlice, item.(T))
	}

	return newSlice
}
