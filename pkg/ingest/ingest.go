// Package ingest ...
package ingest

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"fmt"
	"os"
)

// <---------------------------------------------------------------------------------------------------->


func ReadFileSystem() map[string]interface{} {
	const jsonPath = "./bin/fileSystem.json"

	fileSystemData, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Reading fileSystem:", err)
		return nil
	}

	var fileSystemMap = map[string]interface{}{}

	err = json.Unmarshal(fileSystemData, &fileSystemMap)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return fileSystemMap
}