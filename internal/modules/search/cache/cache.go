// Package cache handles everything that has to do with the generation of the cache for the Search function, to the generation of our folder structure and importing of the config.
package cache

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/skillptm/Bolt/internal/config"
	"github.com/skillptm/Bolt/internal/util"
)

// Filesystem stores some metadata for our searches, aswell as the cache of files on the system
type Filesystem struct {
	DefaultDirs            Dirs
	ExcludedDirs           DirsRules
	ExcludeFromDefaultDirs DirsRules
	ExtendedDirs           Dirs
	MaxCPUThreads          int
}

/*
Dirs store the DirMap, which is structuered for us the be able to search through as fast as possible.
They also store Paths, which is a map with which you can access the path to the file. This exists to save memory, for not having to store the same path several times with the file directly.

paths: map[unique ID]Absolute Path

dirMap: map[File Extension]map[File Length][]File{encodedName, Name, pathKey}
*/
type Dirs struct {
	DirMap   map[string]map[int][]File
	BaseDirs map[string]bool
	Mu       sync.Mutex
	Paths    map[int]string
}

// File stores all the data we need for a fast retrival later on
type File struct {
	EncodedName [8]byte
	Name        string
	PathKey     int
}

// DirsRules holds name, path and regex rules determening the part of the cache a folder will be in
type DirsRules struct {
	Name  map[string]bool
	Path  map[string]bool
	Regex []string
}

// basicFile is a temp struct we use to not have to re-gather file data between different actions
type basicFile struct {
	extension string
	isFolder  bool
	name      string
	path      string
}

// NewFilesystem returns a pointer to a Filesystem struct that has been filled up according to the includedDirs, excludedDirs and config
func NewFilesystem() (*Filesystem, error) {

	config, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("NewFilesystem: couldn't create config:\n--> %w", err)
	}

	fs := Filesystem{
		DefaultDirs: Dirs{
			Paths:    make(map[int]string),
			BaseDirs: util.MakeBoolMap(config.DefaultDirs),
			DirMap:   make(map[string]map[int][]File),
		},
		ExcludedDirs: DirsRules{
			util.MakeBoolMap(config.ExcludeDirs["Name"]),
			util.MakeBoolMap(config.ExcludeDirs["Path"]),
			config.ExcludeDirs["Regex"],
		},
		ExcludeFromDefaultDirs: DirsRules{
			util.MakeBoolMap(config.ExcludeFromDefaultDirs["Name"]),
			util.MakeBoolMap(config.ExcludeFromDefaultDirs["Path"]),
			config.ExcludeFromDefaultDirs["Regex"],
		},
		ExtendedDirs: Dirs{
			Paths:    make(map[int]string),
			BaseDirs: util.MakeBoolMap(config.ExtendedDirs),
			DirMap:   make(map[string]map[int][]File),
		},
		MaxCPUThreads: config.MaxCPUThreads,
	}

	fs.Update(&fs.DefaultDirs, &fs.ExtendedDirs)
	fs.Update(&fs.ExtendedDirs, &fs.DefaultDirs)

	go fs.autoUpdateCache(config.DefaultDirsCacheUpdateTime, config.ExtendedDirsCacheUpdateTime)

	return &fs, nil
}

// Check finds out if the provided Directory breaks any of the name, path or regex rules
func (dr *DirsRules) Check(dirPath string, add bool, dirs *Dirs) (bool, error) {
	addPath := func() {
		if !add {
			return
		}
		dirs.Mu.Lock()
		dirs.BaseDirs[dirPath] = true
		dirs.Mu.Unlock()
	}

	if dr.Path[dirPath] {
		addPath()
		return false, nil
	}

	if dr.Name[path.Base(dirPath)] {
		addPath()
		return false, nil
	}

	for _, pattern := range dr.Regex {
		if matched, err := regexp.MatchString(pattern, dirPath); matched {
			addPath()
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("Check: couldn't match pattern %s:\n--> %w", pattern, err)
		}
	}

	return true, nil
}

// autoUpdateCache automatically updates both the DefaultDirs and ExtendedDirs
func (fs *Filesystem) autoUpdateCache(defaultTime int, extendedTime int) {
	defaultTimer := time.NewTimer(time.Duration(defaultTime) * time.Second)
	extendedTimer := time.NewTimer(time.Duration(extendedTime) * time.Second)

	for {
		select {
		case <-defaultTimer.C:
			fs.Update(&fs.DefaultDirs, &fs.ExtendedDirs)
			defaultTimer.Reset(time.Duration(defaultTime) * time.Second)
		case <-extendedTimer.C:
			fs.Update(&fs.ExtendedDirs, &fs.DefaultDirs)
			extendedTimer.Reset(time.Duration(extendedTime) * time.Second)
		}
	}
}

// Update launches the traversing of the dirs and later starts the adding of the results onto the fs
func (fs *Filesystem) Update(dirs *Dirs, otherDirs *Dirs) {

	// 10000000 is the channel size, because we just need a ridiculously large channel to store all the paths until we traversed them
	pathQueue := make(chan string, 10000000)
	results := make(chan *basicFile, 10000000)
	wg := sync.WaitGroup{}

	for dir := range dirs.BaseDirs {
		wg.Add(1)
		pathQueue <- dir
	}

	for range fs.MaxCPUThreads {
		go fs.traverse(pathQueue, results, otherDirs, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
		close(pathQueue)
	}()

	dirs.add(results)
}

// traverse walks through and expands the pathQueue to store all files and folders it encounters in resultsChan unless it breaks with excludedDirs
func (fs *Filesystem) traverse(pathQueue chan string, results chan<- *basicFile, otherDirs *Dirs, wg *sync.WaitGroup) {
	for currentDir := range pathQueue {
		currentEntries, err := os.ReadDir(currentDir)
		// an error here simply means we didn't have the permissions to read a dir, so we ignore it
		if err != nil {
			wg.Done()
			continue
		}

		for _, entry := range currentEntries {
			if entry.IsDir() {
				if entry.Name() == "." || entry.Name() == ".." {
					continue
				}

				entryPath := fmt.Sprintf("%s%s", filepath.Join(currentDir, entry.Name()), string(filepath.Separator))

				if checked, err := fs.ExcludedDirs.Check(entryPath, false, &fs.ExtendedDirs); !checked || err != nil {
					continue
				}

				if checked, err := fs.ExcludeFromDefaultDirs.Check(entryPath, true, &fs.ExtendedDirs); !checked || err != nil {
					continue
				}

				if _, ok := otherDirs.BaseDirs[entryPath]; ok {
					continue
				}

				results <- &basicFile{"Folder", true, entry.Name(), entryPath}
				wg.Add(1)
				pathQueue <- entryPath
			} else {
				fileExtension := filepath.Ext(filepath.Join(currentDir, entry.Name()))
				fileName, _ := strings.CutSuffix(entry.Name(), fileExtension)

				results <- &basicFile{fileExtension, false, fileName, currentDir}
			}
		}

		wg.Done()
	}
}

// add formats and overwrites the dirMap on fs
func (dirs *Dirs) add(results <-chan *basicFile) {
	newDirMap := make(map[string]map[int][]File)
	tempPaths := make(map[string]int)

	for item := range results {
		itemExtension := (*item).extension
		itemName := (*item).name
		itemPath := (*item).path

		if _, ok := newDirMap[itemExtension]; !ok {
			newDirMap[itemExtension] = make(map[int][]File)
		}

		if _, ok := newDirMap[itemExtension][len(itemName)]; !ok {
			newDirMap[itemExtension][len(itemName)] = []File{}
		}

		if _, ok := tempPaths[itemPath]; (*item).isFolder && !ok {
			tempPaths[itemPath] = len(tempPaths)
		}

		newDirMap[itemExtension][len(itemName)] = append(newDirMap[itemExtension][len(itemName)], File{Encode(itemName), itemName, tempPaths[itemPath]})
	}

	newPaths := make(map[int]string)
	for key, value := range tempPaths {
		newPaths[value] = key
	}

	dirs.Paths = newPaths
	dirs.DirMap = newDirMap
}
