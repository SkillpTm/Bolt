// Package config handles the generation of a new config with the modules default values.
package config

// <---------------------------------------------------------------------------------------------------->

import (
	"fmt"
	"math"
	"runtime"

	"github.com/skillptm/Bolt/internal/modules/search/util"
)

// <---------------------------------------------------------------------------------------------------->

var BWSConfig *Config

var DefaultConfig = map[string]interface{}{
	"cpuThreads": int(math.Ceil(float64(runtime.NumCPU()) / float64(4))),
	"mainDirs": []string{
		"C:/Users/<USERNAME>/",
	},
	"excludeSubMainDirs": []string{
		"C:/Users/<USERNAME>/AppData/Roaming",
	},
	"secondaryDirs": []string{
		"C:/",
	},
	"excludeDirs": []string{
		"C:/Windows/",
		"C:/$Recycle.Bin/",
		"C:/Users/<USERNAME>/AppData/Local",
		"C:/Users/<USERNAME>/AppData/LocalLow",
	},
	"excludeDirsByName": []string{
		".git",
		"bin",
		"node_modules",
		"steamapps",
	},
}

// <---------------------------------------------------------------------------------------------------->

type Config struct {
	CPUThreads         int
	MainDirs           []string
	ExcludeSubMainDirs []string
	SecondaryDirs      []string
	ExcludeDirs        []string
	ExcludeDirsByName  []string
}

// <---------------------------------------------------------------------------------------------------->

// New creates a new Conifg struct with the values from ./configs/config.json
func New(configMap map[string]interface{}) (*Config, error) {
	newConfig := Config{}

	newConfig.CPUThreads = configMap["cpuThreads"].(int)
	delete(configMap, "cpuThreads")

	// populate the newConfig with properly formated paths
	for key, value := range configMap {
		newSlice := value.([]string)

		for index, element := range newSlice {
			newSlice[index] = util.FormatEntry(element, true)
		}
		newSlice, err := util.InsertUsername(newSlice)
		if err != nil {
			return &newConfig, fmt.Errorf("couldn't replace '<USERNAME>'; %s", err.Error())
		}

		switch key {
		case "mainDirs":
			newConfig.MainDirs = newSlice
		case "excludeSubMainDirs":
			newConfig.ExcludeSubMainDirs = newSlice
		case "secondaryDirs":
			newConfig.SecondaryDirs = newSlice
		case "excludeDirs":
			newConfig.ExcludeDirs = newSlice
		case "excludeDirsByName":
			newConfig.ExcludeDirsByName = newSlice
		}
	}

	return &newConfig, nil
}
