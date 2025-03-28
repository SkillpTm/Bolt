// Package ssljson includes functions that make the use of JSONs easier for me
package ssljson

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// <---------------------------------------------------------------------------------------------------->

// GetData will provide a map with JSON data of the prvoided file
func GetData(filePath string) (map[string]interface{}, error) {
	var jsonData map[string]interface{}

	// open JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't open JSON file (%s); %s", filePath, err.Error())
	}

	// defer close the file with error handling
	var returnErr error = nil
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("couldn't close JSON file (%s); %s", filePath, jsonFile.Close())
		}
	}()

	// read JSON data from file
	rawJSONData, err := io.ReadAll(jsonFile)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't read JSON file (%s); %s", filePath, err.Error())
	}

	// convert JSON data to map
	err = json.Unmarshal(rawJSONData, &jsonData)
	if err != nil {
		return jsonData, fmt.Errorf("couldn't unmarshal raw JSON data; %s", err.Error())
	}

	return jsonData, returnErr
}

// OverwriteData will take a map with JSON data and the file path and overwrite the data in the existing file
func OverwriteData(filePath string, inputData map[string]interface{}) error {
	// open JSON file in write-only and truncate mode
	jsonFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open JSON file (%s); %s", filePath, err.Error())
	}

	// defer close the file with error handling
	var returnErr error = nil
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			returnErr = fmt.Errorf("couldn't close JSON file (%s); %s", filePath, jsonFile.Close())
		}
	}()

	// marshal the map into JSON
	jsonData, err := json.MarshalIndent(inputData, "", "	")
	if err != nil {
		return fmt.Errorf("couldn't marshal JSON data; %s", err.Error())
	}

	// write jsonData to file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("couldn't write JSON data to JSON file (%s); %s", filePath, err.Error())
	}

	return returnErr
}
