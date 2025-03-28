// Package search handles the search through the cache and the ranking of the results.
package search

import (
	"io/fs"
	"math"
	"os"
	"sync"
	"time"

	"github.com/skillptm/bws/internal/config"
)

const (
	minimumFileSize    int64 = 100 // in bytes
	fourYearsInSeconds int64 = 4 * 365.25 * 24 * 60 * 60

	exactMatchModifier    int     = 500
	minimumSizeModifier   int     = 25
	timeSinceMaxModifier  float64 = 200
	nameLengthMaxModifier float64 = 100
)

// RankedFile holds the points given to a file and it's full path
type RankedFile struct {
	path   string
	points int
}

// newRankedFile constructs a RankedFile and ranks it based on: exact match, minimum file size,
func newRankedFile(fileInfo fs.FileInfo, file []string, pattern *SearchString) *RankedFile {
	newFile := RankedFile{path: file[0]}

	// check if the searchString and the file name are an exact match (except for case)
	if file[1] == pattern.name {
		newFile.points += exactMatchModifier
	}

	// check if the size is of a minimum file size
	if fileInfo.Size() > minimumFileSize {
		newFile.points += minimumSizeModifier
	}

	timeSinceMod := time.Now().UTC().Unix() - fileInfo.ModTime().UTC().Unix()

	// rank how long ago the file was last modified (longer ago = worse)
	if timeSinceMod > fourYearsInSeconds {
		newFile.points += 0
	} else {

		timeSinceReduction := 1 - math.Round(float64(timeSinceMod)/float64(fourYearsInSeconds)*math.Pow(10, 2))/math.Pow(10, 2)

		newFile.points += int(timeSinceMaxModifier * timeSinceReduction)
	}

	// rank how long the filename is compared to the searchString (longer = worse)
	nameLengthReduction := math.Round(float64(pattern.length)/float64(len(file[1]))*math.Pow(10, 2)) / math.Pow(10, 2)
	newFile.points += int(nameLengthMaxModifier * nameLengthReduction)

	return &newFile
}

// Rank ranks and sorts the results
func Rank(searchResults *[][]string, pattern *SearchString, forceStopChan chan bool) *[]string {
	output := []string{}
	rankedFiles := []RankedFile{}

	if len(*searchResults) < 1 {
		return &output
	}

	var wg sync.WaitGroup

	// 10000000 is the channel size, because we just need a ridiculously large channel to store all the results
	toRankChan := make(chan *[]string, 10000000)
	rankedChan := make(chan *RankedFile, 10000000)
	breakChan := make(chan bool, config.BWSConfig.CPUThreads*config.BWSConfig.CPUThreads)

	// we add all the results in the channel before starting the goroutines to avoid race conditions
	for _, file := range *searchResults {
		toRankChan <- &file
	}

	for range config.BWSConfig.CPUThreads {
		wg.Add(1)
		go rankResults(toRankChan, rankedChan, breakChan, pattern, &wg, forceStopChan)
	}

	wg.Wait()

	close(toRankChan)
	close(rankedChan)
	close(breakChan)

	// check if we have to stop the baseSearch
	if len(forceStopChan) > 0 {
		return &output
	}

	for file := range rankedChan {
		rankedFiles = append(rankedFiles, *file)
	}

	// sort the results
	quickSort(rankedFiles)

	// put the ranked and sorted paths onto the output
	for _, file := range rankedFiles {
		output = append(output, file.path)
	}

	return &output
}

// rankResults takes the results from toRankChan, ranks them and inserts a pointer to them into rankedChan
func rankResults(toRankChan <-chan *[]string, rankedChan chan<- *RankedFile, breakChan chan bool, pattern *SearchString, wg *sync.WaitGroup, forceStopChan chan bool) {
	defer wg.Done()

	for {
		select {
		case file := <-toRankChan:
			// check if we have to stop the baseSearch
			if len(forceStopChan) > 0 {
				return
			}

			fileInfo, err := os.Stat((*file)[0])
			if err != nil {
				// if we error it's most likely the file doesn't exist anymore, so we skip it
				continue
			}

			rankedChan <- newRankedFile(fileInfo, *file, pattern)

			if len(toRankChan) < 1 {
				for range config.BWSConfig.CPUThreads {
					breakChan <- true
				}
				return
			}

		case <-breakChan:
			return
		}
	}
}

// quickSort is an implmentation of the quick sort alogirthm that sorts our ranked files based on their points
func quickSort(rankedFiles []RankedFile) {
	if len(rankedFiles) <= 1 {
		return
	}

	pivotIndex := len(rankedFiles) / 2
	pivot := rankedFiles[pivotIndex].points

	// partition the slice into two halves
	left := 0
	right := len(rankedFiles) - 1

	for left <= right {
		for rankedFiles[left].points > pivot {
			left++
		}

		for rankedFiles[right].points < pivot {
			right--
		}

		if left <= right {
			rankedFiles[left], rankedFiles[right] = rankedFiles[right], rankedFiles[left]
			left++
			right--
		}
	}

	// recursively sort the two partitions
	quickSort(rankedFiles[:right+1])
	quickSort(rankedFiles[right+1:])
}
