// Package modules holdes the interfaces for core modules of Bolt
package modules

import (
	"fmt"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/skillptm/Bolt/internal/config"
	"github.com/skillptm/Bolt/internal/modules/search"
	"github.com/skillptm/Bolt/internal/modules/search/cache"
	"github.com/skillptm/Bolt/internal/util"
)

// SearchHandler is an interface which will hold the indexed cache and be the start point for searches
type SearchHandler struct {
	fileSystem    *cache.Filesystem
	forceStopChan chan bool
	searching     bool

	ResultsChan chan []string
}

// NewSearchHandler is the constructor for SearchHandler, which also sets up the cache and the Filesystem
func NewSearchHandler(conf *config.Config) (*SearchHandler, error) {
	sh := SearchHandler{
		forceStopChan: make(chan bool, 1),
		ResultsChan:   make(chan []string, 1),
		searching:     false,
	}

	fs, err := cache.NewFilesystem(conf)
	if err != nil {
		return nil, fmt.Errorf("NewFilesystem: couldn't setup Filesystem:\n--> %w", err)
	}

	sh.fileSystem = fs

	return &sh, nil
}

// ClearImportedCache clears the cache data from memory
func (sh *SearchHandler) ClearImportedCache() {
	sh.fileSystem.DefaultDirs.DirMap = make(map[string]map[int][]cache.File)
	sh.fileSystem.DefaultDirs.Paths = make(map[int]string)
	sh.fileSystem.DefaultDirs.Imported = false

	sh.fileSystem.ExtendedDirs.DirMap = make(map[string]map[int][]cache.File)
	sh.fileSystem.ExtendedDirs.Paths = make(map[int]string)
	sh.fileSystem.ExtendedDirs.Imported = false

	runtime.GC()
	debug.FreeOSMemory()
}

// ForceUpdateCache immediately updates the cache. If extended is set the default/extended caches are updated and reset will reset the whole Filesystem on the SearchHandler
func (sh *SearchHandler) ForceUpdateCache(conf *config.Config, extended bool, reset bool) error {
	if reset {
		fs, err := cache.NewFilesystem(conf)
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

// ImportCache imports the cache data from the disk into memory
func (sh *SearchHandler) ImportCache() {
	sh.fileSystem.DefaultDirs.Mu.Lock()
	util.GetJSON(sh.fileSystem.DefaultDirs.CachePath, &sh.fileSystem.DefaultDirs)
	sh.fileSystem.DefaultDirs.Mu.Unlock()
	sh.fileSystem.DefaultDirs.Imported = true

	// in a goroutine to speed up start up time
	go func() {
		sh.fileSystem.ExtendedDirs.Mu.Lock()
		util.GetJSON(sh.fileSystem.ExtendedDirs.CachePath, &sh.fileSystem.ExtendedDirs)
		sh.fileSystem.ExtendedDirs.Mu.Unlock()
		sh.fileSystem.ExtendedDirs.Imported = true
	}()
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

	// Importing the extended dirs is done over a goroutine, which might not have finished here. So we wait for it and break early with the forceStopChan, if needed
	if extendedSearch && !sh.fileSystem.ExtendedDirs.Imported {
		for !sh.fileSystem.ExtendedDirs.Imported {
			time.Sleep(time.Duration(5) * time.Millisecond)
		}

		if len(sh.forceStopChan) > 0 {
			return
		}
	}

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
<file extensions>: which tells us the file extensions. The separator for extensions is a ','

Example:

input: "myFile /e <txt, go>" -> output: "myfile", ["txt", "go"], true
*/
func matchFlags(input string) (string, []string, bool) {
	input = strings.ToLower(input)
	extendedSearch := false
	extensions := []string{}

	// the pattern detects: /e and /E for the extended search flag
	pattern := "/e"

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
			for _, char := range []string{"<", ">", " "} {
				match = strings.ReplaceAll(match, char, "")
			}

			extensions = append(extensions, strings.Split(match, ",")...)
		}
	}

	input = regex.ReplaceAllString(input, "")

	// remove any lone flag characters from the search
	input = strings.Trim(input, " /<>")

	if index := strings.LastIndex(input, "."); index >= 0 && !slices.Contains(extensions, "folder") {
		extensions = append(extensions, input[index:])
		input = input[:index]
	}

	return input, extensions, extendedSearch
}
