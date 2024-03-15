// Package searchhandler ...
package searchhandler

// <---------------------------------------------------------------------------------------------------->

import (
	"regexp"
	"strings"

	"github.com/skillptm/bws"
)

// <---------------------------------------------------------------------------------------------------->

// SearchHandler is in charge of searching and stoping searches
type SearchHandler struct {
	BreakChan   chan bool
	ResultsChan chan []string

	searching bool
}

// New returns a new SearchHandler
func New() *SearchHandler {
	return &SearchHandler{
		make(chan bool, 1),
		make(chan []string, 1),
		false,
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

	// the pattern removes: /e and /E for the extended search flag
	pattern := "(/e|/E)"

	regex := regexp.MustCompile(pattern)

	// if we found anything the extended search flag was set
	if len(regex.FindAllString(input, -1)) > 0 {
		extendedSearch = true
	}

	input = regex.ReplaceAllString(input, "")

	// the pattern removes: anything between (and including) < and > for the extensions
	pattern = "<[^>]*>"

	regex = regexp.MustCompile(pattern)

	if matches := regex.FindAllString(input, -1); len(matches) > 0 {
		for _, match := range matches {
			// remove the flag chars and spaces
			for _, char := range [3]string{"<", ">", " "} {
				match = strings.ReplaceAll(match, char, "")
			}

			// put the extensions split by the sperator ',' onto the extensions
			extensions = append(extensions, strings.Split(match, ",")...)
		}
	}

	input = regex.ReplaceAllString(input, "")

	// remove any lone flag characters from the search
	for _, char := range [3]string{"/", "<", ">"} {
		input = strings.ReplaceAll(input, char, "")
	}

	input = strings.TrimSpace(input)

	return input, extensions, extendedSearch
}

// StartSearch will start the search in a goroutine
func (searchHandler *SearchHandler) StartSearch(input string) {
	// if there currently is a search ongoing stop it
	if searchHandler.searching {
		searchHandler.BreakChan <- true
		searchHandler.searching = false
	}

	// set a new break channel, because the last one got closed, if we stopped a search
	searchHandler.BreakChan = make(chan bool, 1)

	// launch the search as a gorutine
	go func() {
		searchString, fileExtensions, extendedSearch := matchFlags(input)
		searchHandler.searching = true

		// start the search with the option of breaking it
		result := bws.GoSearchWithBreak(searchString, fileExtensions, extendedSearch, searchHandler.BreakChan)

		// put the results in the results channel so app.EmitSearchResult can emit them to the frontend
		searchHandler.ResultsChan <- result
	}()
}
