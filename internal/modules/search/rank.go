// Package search handles the search, aswell as ranking and sorting of the results.
package search

import (
	"io/fs"
	"strconv"
	"strings"
	"time"
)

const (
	fourYearsInSeconds int   = 4 * 365.25 * 24 * 60 * 60
	minimumSizeAmount  int64 = 100 // in bytes

	exactMatch          int     = 500
	subStringEarlyMax   int     = 325
	recentlyModifiedMax float64 = 250
	notDeeplyNestedMax  int     = 150
	lengthDifferenceMax float64 = 125
	inDefaultDirs       int     = 75
	minimumSize         int     = 25
)

// rankedFile holds the points given to a file and it's full path
type rankedFile struct {
	path   string
	points int
}

// newRankedFile constructor for rankedFile
func newRankedFile(fileInfo fs.FileInfo, file []string, filePath string, pattern *SearchString, defaultDirs map[string]bool) *rankedFile {
	newFile := rankedFile{filePath, 0}

	if file[1] == pattern.name {
		newFile.points += exactMatch
	}

	if index, err := strconv.Atoi(file[3]); err == nil {
		newFile.points += subStringEarlyMax - (10 * index)
	}

	modifiedSecondsAgo := min(time.Now().UTC().Unix()-fileInfo.ModTime().UTC().Unix(), int64(fourYearsInSeconds))
	newFile.points += int(recentlyModifiedMax * (1 - float64(modifiedSecondsAgo)/float64(fourYearsInSeconds)))

	newFile.points += notDeeplyNestedMax + (-10 * strings.Count(file[0], "/"))

	newFile.points += int(lengthDifferenceMax * float64(len(pattern.name)) / float64(len(file[1])))

	for dir := range defaultDirs {
		if strings.HasPrefix(file[0], dir) {
			newFile.points += inDefaultDirs
			break
		}
	}

	if fileInfo.Size() > minimumSizeAmount {
		newFile.points += minimumSize
	}

	return &newFile
}

// quickSort is an implmentation of the quick sort alogirthm that sorts the ranked files based on their points
func quickSort(rankedFiles []rankedFile) {
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
