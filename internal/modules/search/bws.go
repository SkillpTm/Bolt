// Package bws contains the main Search function and start up logic of bws.
package bws

// <---------------------------------------------------------------------------------------------------->

import (
	"log"
	"runtime"
	"time"

	"github.com/skillptm/bws/internal/cache"
	"github.com/skillptm/bws/internal/config"
	"github.com/skillptm/bws/internal/search"
)

// <---------------------------------------------------------------------------------------------------->

// init executes as soon as bws is imported into another project and launches a goroutine of updateCache
func init() {
	var err error
	config.BWSConfig, err = config.New(config.DefaultConfig)
	if err != nil {
		log.Fatal(err)
	}

	go updateCache()
}

// updateCache updates with the use of tickers both the MainDirs and the SecondaryDirs
func updateCache() {
	// create tickers for how often the FileSystem components are supossed to update
	mainDirsTicker := time.NewTicker(3 * time.Minute)
	defer mainDirsTicker.Stop()

	secondaryDirsTicker := time.NewTicker(30 * time.Minute)
	defer secondaryDirsTicker.Stop()

	for {
		select {
		case <-mainDirsTicker.C:
			if !cache.EntrieFilesystem.SetupProperly {
				continue
			}

			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.MainDirs, true)
				runtime.GC()
			}
		case <-secondaryDirsTicker.C:
			if !cache.EntrieFilesystem.SetupProperly {
				continue
			}

			// check if the Filesystem is Updateable, if so update it, otherwise wait for the next cycle
			if cache.EntrieFilesystem.Updateable {
				cache.EntrieFilesystem.Update(config.BWSConfig.SecondaryDirs, false)
				runtime.GC()
			}
		}
	}
}

/*
baseSearch is a wrapper around the search and rank functions and returns the ranked search results at the end.

Additionally if baseSearch was launched by GoSearchWithBreak it can be stopped at any point with the forceStopChan.
If the search was stopped early we return a true, otherwise false.
*/
func baseSearch(searchString string, fileExtensions []string, extendedSearch bool, forceStopChan chan bool) ([]string, bool) {
	// check if the FileSystem is setup properly, if not reset it by regenerating it
	if !cache.EntrieFilesystem.SetupProperly {
		ForceUpdateCache()
	}

	// make it so while we search we can't update the FileSystem
	cache.EntrieFilesystem.Updateable = false
	defer func() {
		cache.EntrieFilesystem.Updateable = true
	}()

	// get the filepaths and names
	results, pattern := search.Start(search.NewSearchString(searchString, fileExtensions), extendedSearch, forceStopChan)

	// check if we have to stop the baseSearch
	if len(forceStopChan) > 0 {
		return []string{}, true
	}

	output := *search.Rank(results, pattern, forceStopChan)

	// check if we have to stop the baseSearch
	if len(forceStopChan) > 0 {
		return []string{}, true
	}

	// rank and sort the files
	return output, false
}

/*
Search takes in any substring that you want to search for through all filenames. You may add any amount of file extensions as well.
The extendedSearch flag dictates, if we search through the SecondaryDirs.
To change the folders included/excluded in the search use the pkg/options set functions.

On it's first execution the function will take longer, as it needs to generate the cache first.
*/
func Search(searchString string, fileExtensions []string, extendedSearch bool) []string {
	results, _ := baseSearch(searchString, fileExtensions, extendedSearch, make(chan bool, 1)) // insert a dummy channel, as it's not needed here
	return results
}

/*
GoSearchWithBreak behaves exactly like Search, the only difference is, it requires a break channel as an input.

This function should be started as a goroutine and it can be cancelled early by sending something in the breakChan.
If it breaks early, it returns a true, otherwise false.
*/
func GoSearchWithBreak(searchString string, fileExtensions []string, extendedSearch bool, breakChan chan bool) ([]string, bool) {
	forceStopChan := make(chan bool, 1)

	go func() {
		// check for when we receive the break signal
		for range breakChan {
			// send something into the foreStopChan to stop the baseSearch
			forceStopChan <- true
		}
	}()

	return baseSearch(searchString, fileExtensions, extendedSearch, forceStopChan)
}

/*
ForceUpdateCache updates the cache regardless of it's state.

This function is generally not needed. Though it can be useful, if you want to generate the cache early, before your first search.
*/
func ForceUpdateCache() {
	cache.EntrieFilesystem = cache.New(config.BWSConfig.MainDirs, config.BWSConfig.SecondaryDirs)
	runtime.GC()
}
