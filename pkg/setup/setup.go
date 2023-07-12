// Package setup ...
package setup

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// <---------------------------------------------------------------------------------------------------->

var fileSystem = map[string]interface{}{}

// <---------------------------------------------------------------------------------------------------->

func Setup(path []string) {
	recReadPath(path)
	createCache()
}

// <---------------------------------------------------------------------------------------------------->

func recReadPath(path []string) {
	entries, err := os.ReadDir(filepath.Join(path...))
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	currentMap := fileSystem

	for _, folder := range path {
		if nextMap, ok := currentMap[folder].(map[string]interface{}); ok {
			currentMap = nextMap
		} else {
			currentMap[folder] = map[string]interface{}{}
			currentMap = currentMap[folder].(map[string]interface{})
		}
	}

	for _, entry := range entries {

		if entry.IsDir() {
			path = append(path, entry.Name())

			recReadPath(path)

			path = path[:len(path)-1]

		} else {
			filePath := filepath.Join(path...) + "\\" + entry.Name()
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Println("Error getting file info:", err)
				continue
			}

			currentMap[entry.Name()] = fileInfo.Size()
		}
	}
}



func createCache() {
	jsonData, err := json.MarshalIndent(fileSystem, "", "	")
	if err != nil {
		fmt.Println("Error generating json data:", err)
		return
	}

	file, err := os.Create("./bin/fileSystem.json")
	if err != nil {
		fmt.Println("Error creating json file:", err)
		return
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing json file:", err)
		return
	}
}