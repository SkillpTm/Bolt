// Package modules holdes the interfaces for core modules of Bolt
package modules

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/skillptm/Bolt/internal/modules/search"
	"github.com/skillptm/Bolt/internal/modules/search/cache"
)

// SearchHandler is an interface which will hold the indexed cache and be the start point for searches
type SearchHandler struct {
	fileSystem    *cache.Filesystem
	forceStopChan chan bool
	searching     bool

	ResultsChan chan []string
}

// NewSearchHandler is the constructor for SearchHandler, which also sets up the cache and the Filesystem
func NewSearchHandler() (*SearchHandler, error) {
	sh := SearchHandler{
		forceStopChan: make(chan bool, 1),
		ResultsChan:   make(chan []string, 1),
		searching:     false,
	}

	fs, err := cache.NewFilesystem()
	if err != nil {
		return nil, fmt.Errorf("NewFilesystem: couldn't setup Filesystem:\n--> %w", err)
	}

	sh.fileSystem = fs

	return &sh, nil
}

// ForceUpdateCache immediately updates the cache. If extended is set the default/extended caches are updated and reset will reset the whole Filesystem on the SearchHandler
func (sh *SearchHandler) ForceUpdateCache(extended bool, reset bool) error {
	if reset {
		fs, err := cache.NewFilesystem()
		if err != nil {
			return fmt.Errorf("ForceUpdateCache: couldn't setup Filesystem:\n--> %w", err)
		}

		sh.fileSystem = fs
	} else if extended {
		sh.fileSystem.Update(&sh.fileSystem.DefaultDirs, &sh.fileSystem.ExtendedDirs)
		sh.fileSystem.Update(&sh.fileSystem.ExtendedDirs, &sh.fileSystem.DefaultDirs)
	} else {
		sh.fileSystem.Update(&sh.fileSystem.DefaultDirs, &sh.fileSystem.ExtendedDirs)
	}

	runtime.GC()

	return nil
}

// Search is the public facing wrapper for the search function, handling breaking old searches and starting new ones
func (sh *SearchHandler) Search(input string) {
	if sh.searching {
		sh.forceStopChan <- true
		sh.searching = false
	}

	// set a new forceStopChan everytime, to stop confusion on what search to break
	sh.forceStopChan = make(chan bool, 1)
	searchString, fileExtensions, extendedSearch := matchFlags(input)
	sh.searching = true

	result := search.Start(searchString, fileExtensions, extendedSearch, sh.fileSystem, sh.forceStopChan)

	// we only want to emit the results, if we got any and we have a search String to avoid updating to no results in the middle of typing
	if len(searchString) > 0 {
		sh.ResultsChan <- result
	}
}

/*
matchFlags cleans the input and returns the flag values in it, it also removes leading and trailing white space.

The flags it matches for are:

/e and /E: which tell us if the search is an extended search
<file extensions>: which tells us the file extensions. The seperator for extensions is a ','

Example:

input: "myFile /e <txt, go>" -> output: "myfile", ["txt", "go"], true
*/
func matchFlags(input string) (string, []string, bool) {
	extendedSearch := false
	extensions := []string{}

	// the pattern detects: /e and /E for the extended search flag
	pattern := "(/e|/E)"

	regex := regexp.MustCompile(pattern)

	if len(regex.FindAllString(input, 1)) > 0 {
		extendedSearch = true
	}

	input = regex.ReplaceAllString(input, "")

	// the pattern detects: anything between (and including) < and > for the extensions
	pattern = "<[^>]*>"

	regex = regexp.MustCompile(pattern)

	if matches := regex.FindAllString(input, -1); len(matches) > 0 {
		for _, match := range matches {
			for _, char := range [3]string{"<", ">", " "} {
				match = strings.ReplaceAll(match, char, "")
			}

			extensions = append(extensions, strings.Split(match, ",")...)
		}
	}

	input = strings.TrimSpace(regex.ReplaceAllString(input, ""))

	return input, extensions, extendedSearch
}
