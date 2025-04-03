// Package config validates our local folder structure and to retrieve the config data
package config

import (
	"fmt"
	"math"
	"runtime"

	"github.com/skillptm/Bolt/internal/util"
)

// Config is made to structure and order the data for the config.json
type Config struct {
	MaxCPUThreadPercentage      float64  `json:"MaxCPUThreadPercentage"`
	ShortCutEnd                 string   `json:"ShortCutEnd"`
	DefaultDirsCacheUpdateTime  int      `json:"DefaultDirsCacheUpdateTime"`
	ExtendedDirsCacheUpdateTime int      `json:"ExtendedDirsCacheUpdateTime"`
	DefaultDirs                 []string `json:"DefaultDirs"`
	ExtendedDirs                []string `json:"ExtendedDirs"`
	ExcludeFromDefaultDirs      Rules    `json:"ExcludeFromDefaultDirs"`
	ExcludeDirs                 Rules    `json:"ExcludeDirs"`

	MaxCPUThreads int               `json:"-"`
	Paths         map[string]string `json:"-"`
}

// Rules is made to structure and order the data for the config.json
type Rules struct {
	Name  []string `json:"Name"`
	Path  []string `json:"Path"`
	Regex []string `json:"Regex"`
}

// NewConfig is the constructor for Config, it imports the data from the config.json
func NewConfig() (*Config, error) {
	newConfig := Config{Paths: make(map[string]string)}
	var err error

	newConfig.Paths["default_cache.json"], newConfig.Paths["extended_cache.json"], newConfig.Paths["config.json"], err = setup()
	if err != nil {
		return nil, fmt.Errorf("NewConfig: couldn't setup folders:\n--> %w", err)
	}

	err = util.GetJSON(newConfig.Paths["config.json"], &newConfig)
	if err != nil {
		return nil, fmt.Errorf("NewConfig: couldn't get JSON map:\n--> %w", err)
	}

	newConfig.MaxCPUThreads = int(math.Ceil(float64(runtime.NumCPU()) * newConfig.MaxCPUThreadPercentage))

	return &newConfig, nil
}
