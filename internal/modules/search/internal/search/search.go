// Package search handles the search through the cache and the ranking of the results.
package search

// <---------------------------------------------------------------------------------------------------->

import (
	"slices"
	"strings"

	"github.com/skillptm/bws/internal/cache"
)

// <---------------------------------------------------------------------------------------------------->

// SearchString holds all the data releated to the searchString input, so we only have to calculate them once
type SearchString struct {
	encoded    [8]byte
	extensions []string
	length     int
	name       string
}

// NewSearchString returns a pointer to a SearchString struct based on the string input
func NewSearchString(searchString string, fileExtensions []string) *SearchString {
	// make sure all extensions begin with a period, unless it's a "File" or a "Folder"
	for index, element := range fileExtensions {
		if len(element) < 1 {
			// pop empty strings
			fileExtensions = append(fileExtensions[:index], fileExtensions[index+1:]...)
			continue
		}

		// ensure "File"/"Folder" have the right case
		if element == "file" || element == "folder" {
			fileExtensions[index] = "F" + element[1:]
			element = "F" + element[1:]
		}

		if !strings.HasPrefix(element, ".") && element != "File" && element != "Folder" {
			fileExtensions[index] = "." + element
			continue
		}
	}

	return &SearchString{
		encoded:    cache.Encode(searchString),
		extensions: fileExtensions,
		length:     len(searchString),
		name:       strings.ToLower(searchString),
	}
}

// Start wraps around the walkFS function and returns all the results from the MainDirs and SecondaryDirs
func Start(pattern *SearchString, extendedSearch bool, forceStopChan chan bool) (*[][]string, *SearchString) {
	output := [][]string{}

	// check the MainDirs for the search string
	output = append(output, *pattern.searchFS(&cache.EntrieFilesystem.MainDirs, forceStopChan)...)

	// check the SecondaryDirs for the search string
	if extendedSearch {
		output = append(output, *pattern.searchFS(&cache.EntrieFilesystem.SecondaryDirs, forceStopChan)...)
	}

	return &output, pattern
}

// searchFS searches one of the provided FileSystem maps, while skiping files for wrong extensions and ecoded values
func (searchString *SearchString) searchFS(dirs *map[string]map[int][][]interface{}, forceStopChan chan bool) *[][]string {
	output := [][]string{}

	// loop over the extensions
	for extension, lengthMaps := range *dirs {
		// check if extensions were provided and if so, if the current extension is a provided one
		if len(searchString.extensions) > 0 && !slices.Contains(searchString.extensions, extension) {
			continue
		}

		// loop over the filename lengths
		for length, fileSlices := range lengthMaps {
			// check if the filename is longer than the searchString
			if length < searchString.length {
				continue
			}

			// loop over the actual files
			for _, file := range fileSlices {
				// check if we have to stop the baseSearch
				if len(forceStopChan) > 0 {
					return &output
				}

				// check if all required letters are inside the filename
				if !cache.CompareBytes(searchString.encoded, file[2].([8]byte)) {
					continue
				}

				// do a substring search over the filename
				if !strings.Contains(file[1].(string), searchString.name) {
					continue
				}

				// if the searchString is inside the filename add it's path and name to the output
				output = append(output, []string{file[0].(string), file[1].(string)})
			}
		}
	}

	return &output
}
