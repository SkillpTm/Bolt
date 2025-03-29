// Package setup ...
package setup

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/skillptm/Bolt/internal/util"
)

/*
Setup validtes all files/folders we need exist

Returns:

	error: an error, if it fails to validate our folder/files.
*/
func Setup() error {
	dirs := map[string]string{
		"~./.config/bolt/config.json":     "~./.config/bolt/",
		"/var/lib/bolt/search-cache.json": "/var/lib/bolt/",
	}

	err := validateFolders(slices.Collect(maps.Values(dirs)))
	if err != nil {
		return fmt.Errorf("Setup: couldn't validate default folders:\n--> %w", err)
	}

	err = validateFiles(slices.Collect(maps.Keys(dirs)))
	if err != nil {
		return fmt.Errorf("Setup: couldn't validate default files:\n--> %w", err)
	}

	return nil
}

/*
validateFolders checks, if our fodler in ~/.config/ and /var/lib/ exist

Parameters:

	dirsToCheck: slice of absolute paths to directories.

Returns:

	error: an error, if it fails to create a folder.
*/
func validateFolders(dirsToCheck []string) error {
	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			continue
		}

		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("validateFolders: failed to create directory at %s:\n--> %w", dir, err)
		}
	}

	return nil
}

/*
validateFiles checks, if config.json and search-cache.json exist

Parameters:

	filesToCheck: slice of absolute paths to files.

Returns:

	error: an error, if it fails to create a file.
*/
func validateFiles(filesToCheck []string) error {
	for _, file := range filesToCheck {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		_, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("validateFiles: failed to create file at %s:\n--> %w", file, err)
		}

		if strings.HasSuffix(file, "config.json") {
			resetConfig()
		}
	}

	return nil
}

/*
resetConfig resets the config file to the default settings

Returns:

	error: an error, if it fails reset the config file.
*/
func resetConfig() error {
	defaultConfig := map[string]any{
		"maxCPUThreadPercentage": 20, // percentage of threads that may be used, always rounding the threads up
		"defaultDirs": map[string][]string{
			"name": {},
			"path": {
				"~/",
			},
			"regex": {},
		},
		"extendedDirs": map[string][]string{
			"name": {},
			"path": {
				"/",
			},
			"regex": {},
		},
		"excludeFromDefaultDirs": map[string][]string{
			"name": {},
			"path": {},
			"regex": {
				`^\..+`,
			},
		},
		"excludeDirs": map[string][]string{
			"name": {
				".git",
				"node_modules",
				"steamapps",
			},
			"path":  {},
			"regex": {},
		},
	}

	err := util.OverwriteJSON("~./.config/bolt/config.json", defaultConfig)
	if err != nil {
		return fmt.Errorf("resetConfig: couldn't reset default config:\n--> %w", err)
	}

	return nil
}
