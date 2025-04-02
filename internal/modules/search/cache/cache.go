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
	DefaultDirs  Dirs
	ExtendedDirs Dirs

	excludedDirs           dirsRules
	excludeFromDefaultDirs dirsRules
	maxCPUThreads          int
}

/*
Dirs store the DirMap, which is structuered for us the be able to search through as fast as possible.
They also store Paths, which is a map with which you can access the path to the file. This exists to save memory, for not having to store the same path several times with the file directly.

paths: map[unique ID]Absolute Path

dirMap: map[File Extension]map[File Length][]File{encodedName, Name, pathKey}
*/
type Dirs struct {
	mu sync.Mutex

	DirMap   map[string]map[int][]File
	BaseDirs map[string]bool
	Paths    map[int]string
}

// File stores all the data we need for a fast retrival later on
type File struct {
	EncodedName [8]byte
	Name        string
	PathKey     int
}

// dirsRules holds name, path and regex rules determening the part of the cache a folder will be in
type dirsRules struct {
	name  map[string]bool
	path  map[string]bool
	regex []string
}

// basicFile is a temp struct we use to not have to re-gather file data between different actions
type basicFile struct {
	extension string
	isFolder  bool
	name      string
	path      string
}

// NewFilesystem returns a pointer to a Filesystem struct that has been filled up according to the includedDirs, excludedDirs and config
func NewFilesystem(conf *config.Config) (*Filesystem, error) {
	fs := Filesystem{
		DefaultDirs: Dirs{
			Paths:    make(map[int]string),
			BaseDirs: util.MakeBoolMap(conf.DefaultDirs),
			DirMap:   make(map[string]map[int][]File),
		},
		excludedDirs: dirsRules{
			util.MakeBoolMap(conf.ExcludeDirs["Name"]),
			util.MakeBoolMap(conf.ExcludeDirs["Path"]),
			conf.ExcludeDirs["Regex"],
		},
		excludeFromDefaultDirs: dirsRules{
			util.MakeBoolMap(conf.ExcludeFromDefaultDirs["Name"]),
			util.MakeBoolMap(conf.ExcludeFromDefaultDirs["Path"]),
			conf.ExcludeFromDefaultDirs["Regex"],
		},
		ExtendedDirs: Dirs{
			Paths:    make(map[int]string),
			BaseDirs: util.MakeBoolMap(conf.ExtendedDirs),
			DirMap:   make(map[string]map[int][]File),
		},
		maxCPUThreads: conf.MaxCPUThreads,
	}

	fs.Update(&fs.DefaultDirs, &fs.ExtendedDirs)
	fs.Update(&fs.ExtendedDirs, &fs.DefaultDirs)

	go fs.autoUpdateCache(conf.DefaultDirsCacheUpdateTime, conf.ExtendedDirsCacheUpdateTime)

	return &fs, nil
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

	for range fs.maxCPUThreads {
		go fs.traverse(pathQueue, results, otherDirs, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
		close(pathQueue)
	}()

	dirs.add(results)
}

// check finds out if the provided Directory breaks any of the name, path or regex rules
func (dr *dirsRules) check(dirPath string, add bool, dirs *Dirs) (bool, error) {
	addPath := func() {
		if !add {
			return
		}
		dirs.mu.Lock()
		dirs.BaseDirs[dirPath] = true
		dirs.mu.Unlock()
	}

	if dr.path[dirPath] {
		addPath()
		return false, nil
	}

	if dr.name[path.Base(dirPath)] {
		addPath()
		return false, nil
	}

	for _, pattern := range dr.regex {
		if matched, err := regexp.MatchString(pattern, dirPath); matched {
			addPath()
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("check: couldn't match pattern %s:\n--> %w", pattern, err)
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

				if checked, err := fs.excludedDirs.check(entryPath, false, &fs.ExtendedDirs); !checked || err != nil {
					continue
				}

				if checked, err := fs.excludeFromDefaultDirs.check(entryPath, true, &fs.ExtendedDirs); !checked || err != nil {
					continue
				}

				if _, ok := otherDirs.BaseDirs[entryPath]; ok {
					continue
				}

				results <- &basicFile{"folder", true, entry.Name(), entryPath}
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
		itemExtension := strings.ToLower((*item).extension)
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
