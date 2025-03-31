// Package search handles the search, aswell as ranking and sorting of the results.
package search

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/skillptm/Bolt/internal/modules/search/cache"
)

// SearchString holds all the data releated to the searchString input, so we only have to calculate them once
type SearchString struct {
	encoded    [8]byte
	extensions []string
	name       string
}

// NewSearchString returns a pointer to a SearchString struct based on the string input
func NewSearchString(searchString string, fileExtensions []string) *SearchString {
	properExtensions := []string{}

	// make sure all extensions begin with a period, unless it's "Folder"
	for _, element := range fileExtensions {
		if element == "" {
			continue
		}

		// ensure "Folder" has the right case
		if element == "folder" {
			element = "Folder"
		}

		if !strings.HasPrefix(element, ".") && element != "Folder" {
			element = "." + element
		}

		properExtensions = append(properExtensions, element)
	}

	return &SearchString{
		encoded:    cache.Encode(searchString),
		extensions: properExtensions,
		name:       strings.ToLower(searchString),
	}
}

// Start wraps around searchFS and then also sorts and ranks the results. The forceStopChan can search it to end it's search early. This will make it yield no results.
func Start(searchString string, fileExtensions []string, extendedSearch bool, fs *cache.Filesystem, forceStopChan chan bool) []string {
	if len(searchString) < 1 {
		return []string{}
	}

	output := []string{}
	pattern := NewSearchString(searchString, fileExtensions)
	foundFilesChan := make(chan *[]string, 10000000)
	rankedFiles := []rankedFile{}
	wg := sync.WaitGroup{}

	wg.Add(1)
	go pattern.searchFS(&fs.DefaultDirs, foundFilesChan, forceStopChan, &wg)

	if extendedSearch {
		wg.Add(1)
		go pattern.searchFS(&fs.ExtendedDirs, foundFilesChan, forceStopChan, &wg)
	}

	go func() {
		wg.Wait()
		close(foundFilesChan)
	}()

	for foundFile := range foundFilesChan {
		if len(forceStopChan) > 0 {
			return output
		}
		fullPath := ""

		if (*foundFile)[2] != "Folder" {
			fullPath = fmt.Sprintf("%s%s%s", (*foundFile)[0], (*foundFile)[1], (*foundFile)[2])
		} else {
			fullPath = (*foundFile)[0]
		}

		fileInfo, err := os.Stat(fullPath)
		// if we error, it's most likely the file doesn't exist anymore, so we skip it
		if err != nil {
			continue
		}

		rankedFiles = append(rankedFiles, *newRankedFile(fileInfo, *foundFile, fullPath, pattern, fs.DefaultDirs.BaseDirs))
	}

	if len(forceStopChan) > 0 {
		return output
	}

	quickSort(rankedFiles)

	for _, rankedFile := range rankedFiles {
		output = append(output, rankedFile.path)
	}

	return output
}

// searchFS searches one of the provided FileSystem maps, while skiping files for wrong extensions and ecoded values
func (searchString *SearchString) searchFS(dirs *cache.Dirs, foundFilesChan chan<- *[]string, forceStopChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	extensionsToCheck := []string{}

	if len(searchString.extensions) > 0 {
		for _, searchExt := range searchString.extensions {
			if _, ok := dirs.DirMap[searchExt]; ok {
				extensionsToCheck = append(extensionsToCheck, searchExt)
			}
		}
	} else {
		extensionsToCheck = slices.Collect(maps.Keys(dirs.DirMap))
	}

	for _, extension := range extensionsToCheck {
		for length, files := range dirs.DirMap[extension] {
			if length < len(searchString.name) {
				continue
			}

			for _, file := range files {
				if len(forceStopChan) > 0 {
					return
				}

				if !cache.CompareEncoding(searchString.encoded, file.EncodedName) {
					continue
				}

				if !strings.Contains(strings.ToLower(file.Name), searchString.name) {
					continue
				}

				foundFilesChan <- &[]string{dirs.Paths[file.PathKey], file.Name, extension}
			}
		}
	}
}
