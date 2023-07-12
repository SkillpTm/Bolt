package main

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"time"

	"github.com/SkillpTm/WindowsSpotlight/pkg/setup"
)

// <---------------------------------------------------------------------------------------------------->

func main() {
	startTime := time.Now()

	path := []string{"C:/"}

	setup.Setup(path)

	elapsedTime := time.Since(startTime)

	fmt.Println("Program execution time:", elapsedTime)
}