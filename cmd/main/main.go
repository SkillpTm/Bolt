package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"time"

	"github.com/SkillpTm/WindowsSpotlight/pkg/ingest"
	"github.com/SkillpTm/WindowsSpotlight/pkg/setup"
)

// <---------------------------------------------------------------------------------------------------->

var fileSystem = map[string]interface{}{}

// <---------------------------------------------------------------------------------------------------->

func main() {
	startTime := time.Now()

	path := []string{"C:/Users/skill/Uni"}

	setup.Setup(path)

	fileSystem = ingest.ReadFileSystem()

	elapsedTime := time.Since(startTime)

	fmt.Println("Program execution time:", elapsedTime)
}