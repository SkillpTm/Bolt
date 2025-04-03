// Package util provides a variation of functions to be used throughout the project
package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// GetJSON will get the data from a JSON file and put it on a map
func GetJSON(filePath string, dataCarrier any) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("GetJSON: couldn't open JSON file %s:\n--> %w", filePath, err)
	}

	var returnErr error = nil
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("GetJSON: couldn't close JSON file %s:\n--> %w", filePath, err)
		}
	}()

	err = json.NewDecoder(jsonFile).Decode(dataCarrier)
	if err != nil {
		return fmt.Errorf("GetJSON: couldn't decode JSON file %s:\n--> %w", filePath, err)
	}

	return returnErr
}

// OverwriteJSON will take a map or struct with JSON data and the file path and overwrite the data in the existing file
func OverwriteJSON(filePath string, indent bool, data any) error {
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
	if indent {
		encoder.SetIndent("", "	")
	}
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("OverwriteJSON: couldn't encode JSON data:\n--> %w", err)
	}

	return returnErr
}
