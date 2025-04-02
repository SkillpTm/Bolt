// Package search handles the search, aswell as ranking and sorting of the results.
package search

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/skillptm/Bolt/internal/modules/search/cache"
)

// searchString holds all the data releated to the searchString input, so we only have to calculate them once
type searchString struct {
	encoded    [8]byte
	extensions []string
	name       string
}

// NewSearchString returns a pointer to a searchString struct based on the string input
func newSearchString(searchInput string, fileExtensions []string) *searchString {
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

	return &searchString{
		encoded:    cache.Encode(searchInput),
		extensions: properExtensions,
		name:       strings.ToLower(searchInput),
	}
}

// Start wraps around searchFS and then also sorts and ranks the results. The forceStopChan can search it to end it's search early. This will make it yield no results.
func Start(searchInput string, fileExtensions []string, extendedSearch bool, fs *cache.Filesystem, forceStopChan chan bool) []string {
	if len(searchInput) < 1 {
		return []string{}
	}

	output := []string{}
	pattern := newSearchString(searchInput, fileExtensions)
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
func (sStr *searchString) searchFS(dirs *cache.Dirs, foundFilesChan chan<- *[]string, forceStopChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	extensionsToCheck := []string{}

	if len(sStr.extensions) > 0 {
		for _, searchExt := range sStr.extensions {
			if _, ok := dirs.DirMap[searchExt]; ok {
				extensionsToCheck = append(extensionsToCheck, searchExt)
			}
		}
	} else {
		extensionsToCheck = slices.Collect(maps.Keys(dirs.DirMap))
	}

	for _, extension := range extensionsToCheck {
		for length, files := range dirs.DirMap[extension] {
			if length < len(sStr.name) {
				continue
			}

			for _, file := range files {
				if len(forceStopChan) > 0 {
					return
				}

				if !cache.CompareEncoding(sStr.encoded, file.EncodedName) {
					continue
				}

				if index := strings.Index(strings.ToLower(file.Name), sStr.name); index >= 0 {
					foundFilesChan <- &[]string{dirs.Paths[file.PathKey], file.Name, extension, strconv.Itoa(index)}
				}
			}
		}
	}
}
