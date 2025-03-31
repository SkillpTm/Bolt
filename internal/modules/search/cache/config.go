// Package cache handles everything that has to do with the generation of the cache for the Search function, to the generation of our folder structure and importing of the config.
package cache

import (
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/skillptm/Bolt/internal/util"
)

// DirsRules holds name, path and regex rules determening the part of the cache a folder will be in
type DirsRules struct {
	Name  map[string]bool
	Path  map[string]bool
	Regex []string
}

// config holds the data from the config.json
type config struct {
	maxCPUThreads          int
	defaultDirs            map[string]bool
	extendedDirs           map[string]bool
	excludeFromDefaultDirs DirsRules
	excludeDirs            DirsRules
}

// Check finds out if the provided Directory breaks any of the name, path or regex rules
func (dr *DirsRules) Check(dirPath string, add bool, dirs *Dirs) (bool, error) {
	addPath := func() {
		if !add {
			return
		}
		dirs.Mu.Lock()
		dirs.BaseDirs[dirPath] = true
		dirs.Mu.Unlock()
	}

	if dr.Path[dirPath] {
		addPath()
		return false, nil
	}

	if dr.Name[path.Base(dirPath)] {
		addPath()
		return false, nil
	}

	for _, pattern := range dr.Regex {
		if matched, err := regexp.MatchString(pattern, dirPath); matched {
			addPath()
			return false, nil
		} else if err != nil {
			return false, fmt.Errorf("Check: couldn't match pattern %s:\n--> %w", pattern, err)
		}
	}

	return true, nil
}

// newConfig is the constructor for Config, it imports the data from ~./.config/bolt/config.json
func newConfig() (*config, error) {
	newConfig := config{}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return &newConfig, fmt.Errorf("newConfig: couldn't access the user's config dir:\n--> %w", err)
	}

	configMap, err := util.GetJSON(fmt.Sprintf("%s/bolt/config.json", configDir))
	if err != nil {
		return &newConfig, fmt.Errorf("newConfig: couldn't get JSON map:\n--> %w", err)
	}

	getPathsMap := func(value []any) map[string]bool {
		output := make(map[string]bool)

		for _, dir := range value {
			if !strings.HasSuffix(dir.(string), "/") {
				dir = fmt.Sprintf("%s/", dir.(string))
			}

			output[dir.(string)] = true
		}

		return output
	}

	getDirsRules := func(value any) DirsRules {
		rules := DirsRules{
			make(map[string]bool),
			getPathsMap(value.(map[string]any)["path"].([]any)),
			[]string{},
		}

		for _, name := range value.(map[string]any)["name"].([]any) {
			rules.Name[name.(string)] = true
		}

		for _, regex := range value.(map[string]any)["regex"].([]any) {
			rules.Regex = append(rules.Regex, regex.(string))
		}

		return rules
	}

	for key, value := range configMap {
		switch key {
		case "maxCPUThreadPercentage":
			newConfig.maxCPUThreads = int(math.Ceil(float64(runtime.NumCPU()) * configMap["maxCPUThreadPercentage"].(float64)))
		case "defaultDirs":
			newConfig.defaultDirs = getPathsMap(value.([]any))
		case "extendedDirs":
			newConfig.extendedDirs = getPathsMap(value.([]any))
		case "excludeFromDefaultDirs":
			newConfig.excludeFromDefaultDirs = getDirsRules(value)
		case "excludeDirs":
			newConfig.excludeDirs = getDirsRules(value)
		}
	}

	return &newConfig, nil
}
