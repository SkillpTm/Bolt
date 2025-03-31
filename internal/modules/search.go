// Package modules holdes the interfaces for core modules of Bolt
package modules

import (
	"fmt"
	"runtime"

	"github.com/skillptm/Bolt/internal/modules/search"
	"github.com/skillptm/Bolt/internal/modules/search/cache"
)

// SearchModule is an interface which will hold the indexed cache and be the start point for searches
type SearchModule struct {
	fileSystem *cache.Filesystem
}

// NewSearchModule is the constructor for SearchModule, which also sets up the cache and the Filesystem
func NewSearchModule() (*SearchModule, error) {
	sm := SearchModule{}

	fs, err := cache.NewFilesystem()
	if err != nil {
		return &sm, fmt.Errorf("NewFilesystem: couldn't setup Filesystem:\n--> %w", err)
	}

	sm.fileSystem = fs

	return &sm, nil
}

// Search is the public facing wrapper for the search function
func (sm *SearchModule) Search(searchString string, fileExtensions []string, extendedSearch bool, forceStopChan chan bool) []string {
	return search.Start(searchString, fileExtensions, sm.fileSystem, extendedSearch, forceStopChan)
}

// ForceUpdateCache immediately updates the cache. If extended is set the default/extended caches are updated and reset will reset the whole Filesystem on the SearchModule
func (sm *SearchModule) ForceUpdateCache(extended bool, reset bool) error {
	if reset {
		fs, err := cache.NewFilesystem()
		if err != nil {
			return fmt.Errorf("ForceUpdateCache: couldn't setup Filesystem:\n--> %w", err)
		}

		sm.fileSystem = fs
	} else if extended {
		sm.fileSystem.Update(&sm.fileSystem.DefaultDirs, &sm.fileSystem.ExtendedDirs)
		sm.fileSystem.Update(&sm.fileSystem.ExtendedDirs, &sm.fileSystem.DefaultDirs)
	} else {
		sm.fileSystem.Update(&sm.fileSystem.DefaultDirs, &sm.fileSystem.ExtendedDirs)
	}

	runtime.GC()

	return nil
}
