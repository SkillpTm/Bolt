// Package config validates our local folder structure and to retrieve the config data
package config

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/skillptm/Bolt/internal/util"
)

// Config holds the data from the config.json
type Config struct {
	DefaultDirs                 []string
	DefaultDirsCacheUpdateTime  int
	ExcludeDirs                 map[string][]string
	ExcludeFromDefaultDirs      map[string][]string
	ExtendedDirs                []string
	ExtendedDirsCacheUpdateTime int
	MaxCPUThreads               int
	Paths                       map[string]string
}

// NewConfig is the constructor for Config, it imports the data from the config.json
func NewConfig() (*Config, error) {
	newConfig := Config{}
	var err error

	newConfig.Paths["search_cache.json"], newConfig.Paths["config.json"], err = setup()
	if err != nil {
		return nil, fmt.Errorf("NewConfig: couldn't setup folders:\n--> %w", err)
	}

	configMap, err := util.GetJSON(newConfig.Paths["config.json"])
	if err != nil {
		return nil, fmt.Errorf("newCoNewConfignfig: couldn't get JSON map:\n--> %w", err)
	}

	getPaths := func(input []any) []string {
		output := []string{}

		for _, dir := range input {
			if !strings.HasSuffix(dir.(string), "/") {
				dir = fmt.Sprintf("%s/", dir.(string))
			}

			output = append(output, dir.(string))
		}

		return output
	}

	getDirsRules := func(input map[string]any) map[string][]string {
		rules := map[string][]string{
			"Name":  {},
			"Path":  getPaths(input["Path"].([]any)),
			"Regex": {},
		}

		for _, name := range input["Name"].([]any) {
			rules["Name"] = append(rules["Name"], name.(string))
		}

		for _, regex := range input["Regex"].([]any) {
			rules["Regex"] = append(rules["Regex"], regex.(string))
		}

		return rules
	}

	for key, value := range configMap {
		switch key {
		case "MaxCPUThreadPercentage":
			newConfig.MaxCPUThreads = int(math.Ceil(float64(runtime.NumCPU()) * configMap["MaxCPUThreadPercentage"].(float64)))
		case "DefaultDirsCacheUpdateTime":
			newConfig.DefaultDirsCacheUpdateTime = int(configMap["DefaultDirsCacheUpdateTime"].(float64))
		case "ExtendedDirsCacheUpdateTime":
			newConfig.ExtendedDirsCacheUpdateTime = int(configMap["ExtendedDirsCacheUpdateTime"].(float64))
		case "DefaultDirs":
			newConfig.DefaultDirs = getPaths(value.([]any))
		case "ExtendedDirs":
			newConfig.ExtendedDirs = getPaths(value.([]any))
		case "ExcludeFromDefaultDirs":
			newConfig.ExcludeFromDefaultDirs = getDirsRules(value.(map[string]any))
		case "ExcludeDirs":
			newConfig.ExcludeDirs = getDirsRules(value.(map[string]any))
		}
	}

	return &newConfig, nil
}
