package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var fileSystem = map[string]interface{}{}

func main() {
	path := "./"

	rec_read_entry(path, nil)

	for key, value := range fileSystem {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func rec_read_entry(path string, entry fs.DirEntry) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			rec_read_entry(entryPath, entry)

		} else {
			entryInfo, err := os.Stat(entryPath)
			if err != nil {
				log.Fatal(err)
			}

			fileSystem[entryPath] = entryInfo.Size()
		}
	}
}
