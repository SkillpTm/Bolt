package setup

import (
	"fmt"
	"math"
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

	configMap, err := util.GetJSON("~./.config/bolt/config.json")
	if err != nil {
		return &newConfig, fmt.Errorf("New: couldn't get JSON map:\n--> %w", err)
	}

	newConfig.MaxCPUThreads = int(math.Ceil(float64(runtime.NumCPU()) * float64(configMap["maxCPUThreadPercentage"].(int))))

	for key, value := range configMap {
		names := map[string]bool{}
		for _, name := range value.(map[string][]string)["name"] {
			names[name] = true
		}

		paths := map[string]bool{}
		for _, path := range value.(map[string][]string)["path"] {
			paths[path] = true
		}

		rules := DirsRules{
			names,
			paths,
			value.(map[string][]string)["regex"],
		}

		switch key {
		case "defaultDirs":
			newConfig.DefaultDirs = value.([]string)
		case "extendedDirs":
			newConfig.ExtendedDirs = value.([]string)
		case "excludeFromDefaultDirs":
			newConfig.ExcludeFromDefaultDirs = rules
		case "excludeDirs":
			newConfig.ExcludeDirs = rules
		}
	}

	return &newConfig, nil
}
