// Package util provides a variation of functions to be used throughout the project
package util

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
GetJSON will get the data from a JSON file and put it on a map

Parameters:

	filePath: absolute path to the JSON file.

Returns:

	map[string]any: a map of the data from the provided JSON file
	error: an error, if a file operation on the provided JSON fails.
*/
func GetJSON(filePath string) (map[string]any, error) {
	jsonData := map[string]any{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return jsonData, fmt.Errorf("GetJSON: couldn't open JSON file %s:\n--> %w", filePath, err)
	}

	var returnErr error = nil
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("GetJSON: couldn't close JSON file %s:\n--> %w", filePath, err)
		}
	}()

	err = json.NewDecoder(jsonFile).Decode(&jsonData)
	if err != nil {
		return jsonData, fmt.Errorf("GetJSON: couldn't decode JSON file %s:\n--> %w", filePath, err)
	}

	return jsonData, returnErr
}

/*
OverwriteJSON will take a map with JSON data and the file path and overwrite the data in the existing file

Parameters:

	filePath: absolute path to the JSON file.
	data: a map that will be used to fill the json file.

Returns:

	error: an error, if a file operation on the provided JSON fails.
*/
func OverwriteJSON(filePath string, data map[string]any) error {
	jsonFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("OverwriteJSON: couldn't open JSON file %s:\n--> %w", filePath, err)
	}

	var returnErr error = nil
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("OverwriteJSON: couldn't close JSON file %s:\n--> %w", filePath, err)
		}
	}()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "	")
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("OverwriteJSON: couldn't encode JSON data:\n--> %w", err)
	}

	return returnErr
}
