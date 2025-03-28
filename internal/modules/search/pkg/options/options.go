// Package options allows you to set values from the configaration of the cache generation and search.
package options

// <---------------------------------------------------------------------------------------------------->

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/skillptm/bws/internal/cache"
	"github.com/skillptm/bws/internal/config"
	"github.com/skillptm/bws/internal/util"
)

// <---------------------------------------------------------------------------------------------------->

/*
SetCPUThreads allows you to set the maximum amount of threads that will be used during the cache generation.

By default this valus is 1/4 of your CPU's threads, while always rounding up to the next integer.
*/
func SetCPUThreads(threads int) error {
	if threads < 0 {
		return errors.New("you can only set the CPU threads to a minimum of 1")
	}

	if threads > runtime.NumCPU() {
		return fmt.Errorf("you can only set the CPU threads to a maximum of %d", runtime.NumCPU())
	}

	config.BWSConfig.CPUThreads = threads

	return nil
}

// setConfigDirs checks if all provided folders exist and then sets them to the correct attribute of BWSConfig
func setConfigDirs(configType string, newDirs []string) error {
	// properly format the provided paths
	for index, element := range newDirs {
		newDirs[index] = util.FormatEntry(element, true)
	}
	newDirs, err := util.InsertUsername(newDirs)
	if err != nil {
		return fmt.Errorf("couldn't replace '<USERNAME>'; %s", err.Error())
	}

	//check if all dirs prvoided exist and aren't a file
	for _, dir := range newDirs {
		if fileInfo, err := os.Stat(dir); err != nil {
			return fmt.Errorf("%s couldn't be added to %s, because it either can't be accessed or doesn't exist", dir, configType)
		} else {
			if !fileInfo.IsDir() {
				return fmt.Errorf("%s couldn't be added to %s, because it's a file and not a folder", dir, configType)
			}
		}
	}

	switch configType {
	case "MainDirs":
		config.BWSConfig.MainDirs = newDirs
	case "ExcludeSubMainDirs":
		config.BWSConfig.ExcludeSubMainDirs = newDirs
	case "SecondaryDirs":
		config.BWSConfig.SecondaryDirs = newDirs
	case "ExcludeDirs":
		config.BWSConfig.ExcludeDirs = newDirs
	}

	return nil
}

/*
SetMainDirs allows you to set the MainDirs for the config that controls the cache generation.

Using this function will cause the cache to regenerate before the next bws.Search execution.

By default this valus is "C:/Users/<USERNAME>/".
*/
func SetMainDirs(newDirs []string) error {
	if len(newDirs) < 1 {
		return errors.New("you need to set at least one MainDirs folder")
	}

	err := setConfigDirs("MainDirs", newDirs)
	if err != nil {
		return err
	}

	cache.EntrieFilesystem.SetupProperly = false

	return nil
}

/*
SetExcludeSubMainDirs allows you to set the ExcludeSubMainDirs for the config that controls the cache generation.
Meaning that the subfolders, of the MainDirs provided here, will be added to the SecondaryDirs for the extened search.

Using this function will cause the cache to regenerate before the next bws.Search execution.

By default this valus is "C:/Users/<USERNAME>/AppData/Roaming".
*/
func SetExcludeSubMainDirs(newDirs []string) error {
	err := setConfigDirs("ExcludeSubMainDirs", newDirs)
	if err != nil {
		return err
	}

	cache.EntrieFilesystem.SetupProperly = false

	return nil
}

/*
SetSecondaryDirs allows you to set the SecondaryDirs for the config that controls the cache generation.
These folders will only be search through, when setting the extenedSearch flag in the bws.Search function.

Using this function will cause the cache to regenerate before the next bws.Search execution.

By default this valus is "C:/".
*/
func SetSecondaryDirs(newDirs []string) error {
	err := setConfigDirs("SecondaryDirs", newDirs)
	if err != nil {
		return err
	}

	cache.EntrieFilesystem.SetupProperly = false

	return nil
}

/*
SetExcludeDirs allows you to set the ExcludeDirs for the config that controls the cache generation.
These sepcific folders will not be included in the cache generation and search at all.

Using this function will cause the cache to regenerate before the next bws.Search execution.

By default this valus is "C:/Windows/", "C:/$Recycle.Bin/", "C:/Users/<USERNAME>/AppData/Local" and "C:/Users/<USERNAME>/AppData/LocalLow".
*/
func SetExcludeDirs(newDirs []string) error {
	err := setConfigDirs("ExcludeDirs", newDirs)
	if err != nil {
		return err
	}

	cache.EntrieFilesystem.SetupProperly = false

	return nil
}

/*
SetExcludeDirsByName allows you to set the ExcludeDirsByName for the config that controls the cache generation.
Folders with this specific name, no matter where they are stored, will not be included in the cache generation and search at all.

Using this function will cause the cache to regenerate before the next bws.Search execution.

By default this valus is ".git", "bin", "node_modules" and "steamapps".
*/
func SetExcludeDirsByName(newDirs []string) {
	config.BWSConfig.ExcludeDirsByName = newDirs

	cache.EntrieFilesystem.SetupProperly = false
}
