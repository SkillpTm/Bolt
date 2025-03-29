package setup

import (
	"fmt"

	"github.com/skillptm/Bolt/internal/util"
)

/*
DirsRules holds name, path and regex rules determening the part of the cache a folder will be in
*/
type DirsRules struct {
	name  []string
	path  []string
	regex []string
}

/*
Config holds the data from the config.json
*/
type Config struct {
	maxCPUThreadPercentage int
	defaultDirs            DirsRules
	extendedDirs           DirsRules
	excludeFromDefaultDirs DirsRules
	excludeDirs            DirsRules
}

/*
NewConfig is the constructor for Config, it imports the data from ~./.config/bolt/config.json

Returns:

	*Config: pointer to Config with the data from the config.json
	error: an error, if it fails to get the config data from config.json.
*/
func NewConfig() (*Config, error) {
	newConfig := Config{}

	configMap, err := util.GetJSON("~./.config/bolt/config.json")
	if err != nil {
		return &newConfig, fmt.Errorf("New: couldn't get JSON map:\n--> %w", err)
	}

	newConfig.maxCPUThreadPercentage = configMap["maxCPUThreadPercentage"].(int)

	for key, value := range configMap {
		rules := DirsRules{
			value.(map[string][]string)["name"],
			value.(map[string][]string)["path"],
			value.(map[string][]string)["regex"],
		}

		switch key {
		case "defaultDirs":
			newConfig.defaultDirs = rules
		case "extendedDirs":
			newConfig.extendedDirs = rules
		case "excludeFromDefaultDirs":
			newConfig.excludeFromDefaultDirs = rules
		case "excludeDirs":
			newConfig.excludeDirs = rules
		}
	}

	return &newConfig, nil
}
