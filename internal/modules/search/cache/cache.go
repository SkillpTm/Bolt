// Package cache handles everything that has to do with the generation of the cache for the Search function.
package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/skillptm/Bolt/internal/modules/search/setup"
)

/*
Filesystem stores all the searchable files in the following pattern, which is optimsed primarily for speed and secondarily for memory efficiency:

paths: map[unique ID]Absolute Path
dirMap: map[File Extension]map[File Length][]File{encodedName, name, pathKey}
*/
type Filesystem struct {
	IncludedDirs []string
	ExcludedDirs []setup.DirsRules

	Paths  map[int]string
	DirMap map[string]map[int][]File
}

// File stores all the data we need for a fast retrival later on
type File struct {
	EncodedName [8]byte
	Name        string
	PathKey     int
}

// basicFile is a temp struct we use to not have to re-gather file data between different actions
type basicFile struct {
	extension string
	isFolder  bool
	name      string
	nameLen   int
	path      string
}

// NewFilesystem returns a pointer to a Filesystem struct that has been filled up according to the includedDirs, excludedDirs and config
func NewFilesystem(includedDirs []string, excludedDirs []setup.DirsRules, config *setup.Config) *Filesystem {
	fs := Filesystem{
		IncludedDirs: includedDirs,
		ExcludedDirs: excludedDirs,
		Paths:        make(map[int]string),
		DirMap:       make(map[string]map[int][]File),
	}

	fs.Update(config)

	return &fs
}

// Update launches the traversing of the includedDirs and later starts the adding of the results onto the fs
func (fs *Filesystem) Update(config *setup.Config) {

	// 10000000 is the channel size, because we just need a ridiculously large channel to store all the paths until we traversed them
	pathQueue := make(chan string, 10000000)
	results := make(chan *basicFile, 10000000)
	var wg sync.WaitGroup

	for _, dir := range fs.IncludedDirs {
		wg.Add(1)
		pathQueue <- dir
	}

	for range config.MaxCPUThreads {
		go fs.traverse(pathQueue, results, &wg)
	}

	wg.Wait()
	close(results)
	close(pathQueue)
	fs.add(results)
}

// traverse walks through and expands the pathQueue to store all files and folders it encounters in resultsChan unless it breaks with excludedDirs
func (fs *Filesystem) traverse(pathQueue chan string, resultsChan chan<- *basicFile, wg *sync.WaitGroup) {
	for currentDir := range pathQueue {
		currentEntries, err := os.ReadDir(currentDir)
		// an error here simply means we didn't have the permissions to read a dir, so we ignore it
		if err != nil {
			wg.Done()
			continue
		}

	entryLoop:
		for _, entry := range currentEntries {
			if entry.IsDir() {
				if entry.Name() == "." || entry.Name() == ".." {
					continue entryLoop
				}

				entryPath := fmt.Sprintf("%s%s", filepath.Join(currentDir, entry.Name()), string(filepath.Separator))

				for _, rules := range fs.ExcludedDirs {
					if checked, err := rules.Check(entryPath); !checked || err != nil {
						continue entryLoop
					}
				}

				resultsChan <- &basicFile{"Folder", true, entry.Name(), len(entry.Name()), entryPath}
				wg.Add(1)
				pathQueue <- entryPath
			} else {
				entryPath := filepath.Join(currentDir, entry.Name())
				fileExtension := filepath.Ext(entryPath)
				fileName := entry.Name()

				if len(fileExtension) < 1 {
					fileExtension = "File"
				} else {
					fileName = fileName[:len(fileName)-len(fileExtension)]
				}
				resultsChan <- &basicFile{fileExtension, false, fileName, len(fileName), entryPath}
			}
		}

		wg.Done()
	}
}

// add formats and overwrites the dirMap on fs
func (fs *Filesystem) add(results <-chan *basicFile) {
	newDirMap := make(map[string]map[int][]File)
	newPaths := make(map[int]string)

	for {
		item, ok := <-results
		if !ok && len(results) < 1 {
			break
		}

		itemExtension := (*item).extension
		itemIsFolder := (*item).isFolder
		itemName := (*item).name
		itemNameLen := (*item).nameLen
		itemPath := (*item).path

		if _, ok := newDirMap[itemExtension]; !ok {
			newDirMap[itemExtension] = make(map[int][]File)
		}

		if _, ok := newDirMap[itemExtension][itemNameLen]; !ok {
			newDirMap[itemExtension][itemNameLen] = []File{}
		}

		if itemIsFolder {
			newPaths[len(newPaths)] = itemPath
		}

		newDirMap[itemExtension][itemNameLen] = append(newDirMap[itemExtension][itemNameLen], File{Encode(itemName), itemName, len(newPaths) - 1})
	}

	fs.Paths = newPaths
	fs.DirMap = newDirMap
}
