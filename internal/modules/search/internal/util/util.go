// Package util carries smaller utility functions that can be reused over the whole project.
package util

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"os/user"
	"strings"
)

// <---------------------------------------------------------------------------------------------------->

// InsertUsername replaces <USERNAME> in any path to the actual username of the current user
func InsertUsername(pathInputs []string) ([]string, error) {
	// Get the current user for their name
	currentUser, err := user.Current()
	if err != nil {
		return pathInputs, fmt.Errorf("couldn't get current user for username; %s", err.Error())
	}

	username := strings.Split(currentUser.Username, "\\")[1]

	for index, path := range pathInputs {
		pathInputs[index] = strings.ReplaceAll(path, "<USERNAME>", username)
	}

	return pathInputs, nil
}

// FormatEntry replaces "\" with "/" and for folders adds a "/" at the end, if there isn't one already
func FormatEntry(entry string, isFolder bool) string {
	entry = strings.ReplaceAll(entry, "\\", "/")

	// if the inputs are just entries skip them
	if !isFolder {
		return entry
	}

	// check if the last char is a "/" if not append it
	if string(entry[len(entry)-1]) != "/" {
		entry += "/"
	}

	return entry
}
