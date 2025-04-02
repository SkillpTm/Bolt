// Package config validates our local folder structure and to retrieve the config data
package config

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/skillptm/Bolt/internal/util"
)

// Config holds the data from the config.json
type Config struct {
	MaxCPUThreads               int
	DefaultDirsCacheUpdateTime  int
	ExtendedDirsCacheUpdateTime int
	DefaultDirs                 []string
	ExtendedDirs                []string
	ExcludeFromDefaultDirs      map[string][]string
	ExcludeDirs                 map[string][]string
}

// NewConfig is the constructor for Config, it imports the data from the config.json
func NewConfig() (*Config, error) {
	err := setup()
	if err != nil {
		return nil, fmt.Errorf("NewConfig: couldn't setup folders:\n--> %w", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("NewConfig: couldn't access the user's config dir:\n--> %w", err)
	}

	configMap, err := util.GetJSON(fmt.Sprintf("%s/bolt/config.json", configDir))
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

	newConfig := Config{}

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
