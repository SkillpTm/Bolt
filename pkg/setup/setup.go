// Package setup ...
package setup

// <---------------------------------------------------------------------------------------------------->

import (
	"encoding/json"
	"log"
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
		log.Fatal(err)
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
				log.Fatal(err)
			}

			currentMap[entry.Name()] = fileInfo.Size()
		}
	}
}



func createCache() {
	jsonData, err := json.MarshalIndent(fileSystem, "", "	")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("./bin/fileSystem.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}