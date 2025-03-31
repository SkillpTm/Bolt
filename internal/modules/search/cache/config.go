package setup

import (
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"runtime"

	"github.com/skillptm/Bolt/internal/util"
)

// DirsRules holds name, path and regex rules determening the part of the cache a folder will be in
type DirsRules struct {
	name  map[string]bool
	path  map[string]bool
	regex []string
}

// Config holds the data from the config.json
type Config struct {
	MaxCPUThreads          int
	DefaultDirs            []string
	ExtendedDirs           []string
	ExcludeFromDefaultDirs DirsRules
	ExcludeDirs            DirsRules
}

// Check finds out if the provided Directory breaks any of the name, path or regex rules
func (dr *DirsRules) Check(dirPath string) (bool, error) {
	if dr.path[dirPath] {
		return false, nil
	}

	if dr.name[path.Base(dirPath)] {
		return false, nil
	}

	for _, pattern := range dr.regex {
		if matched, err := regexp.MatchString(pattern, path.Base(dirPath)); matched {
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("Check: couldn't match pattern %s:\n--> %w", pattern, err)
		}
	}

	return true, nil
}

// NewConfig is the constructor for Config, it imports the data from ~./.config/bolt/config.json
func NewConfig() (*Config, error) {
	newConfig := Config{}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return &newConfig, fmt.Errorf("NewConfig: couldn't access the user's config dir:\n--> %w", err)
	}

	configMap, err := util.GetJSON(fmt.Sprintf("%s/bolt/config.json", configDir))
	if err != nil {
		return &newConfig, fmt.Errorf("New: couldn't get JSON map:\n--> %w", err)
	}

	newConfig.MaxCPUThreads = int(math.Ceil(float64(runtime.NumCPU()) * configMap["maxCPUThreadPercentage"].(float64)))

	for key, value := range configMap {
		getStrings := func(input any) []string {
			dirs := []string{}
			for _, dir := range input.([]any) {
				dirs = append(dirs, dir.(string))
			}
			return dirs
		}

		getDirsRules := func() DirsRules {
			names := map[string]bool{}
			for _, name := range value.(map[string]any)["name"].([]any) {
				names[name.(string)] = true
			}

			paths := map[string]bool{}
			for _, path := range value.(map[string]any)["path"].([]any) {
				paths[path.(string)] = true
			}

			rules := DirsRules{
				names,
				paths,
				getStrings(value.(map[string]any)["regex"]),
			}

			return rules
		}

		switch key {
		case "defaultDirs":
			newConfig.DefaultDirs = getStrings(value)
		case "extendedDirs":
			newConfig.ExtendedDirs = getStrings(value)
		case "excludeFromDefaultDirs":
			newConfig.ExcludeFromDefaultDirs = getDirsRules()
		case "excludeDirs":
			newConfig.ExcludeDirs = getDirsRules()
		}
	}

	return &newConfig, nil
}
