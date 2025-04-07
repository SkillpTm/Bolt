// Package config validates our local folder structure and to retrieve the config data
package config

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/skillptm/Bolt/internal/util"
)

// Setup validates all files/folders we need to exist and returns their paths
func setup(icon embed.FS) ([]string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("setup: couldn't access the user's cache dir:\n--> %w", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("setup: couldn't access the user's config dir:\n--> %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("setup: couldn't access the user's home dir:\n--> %w", err)
	}

	err = validateFolders([]string{
		fmt.Sprintf("%s/bolt/", cacheDir),
		fmt.Sprintf("%s/bolt/", configDir),
		fmt.Sprintf("%s/.local/share/bolt/", homeDir),
	})
	if err != nil {
		return nil, fmt.Errorf("setup: couldn't validate default folders:\n--> %w", err)
	}

	files := []string{
		fmt.Sprintf("%s/bolt/default_cache.json", cacheDir),
		fmt.Sprintf("%s/bolt/extended_cache.json", cacheDir),
		fmt.Sprintf("%s/bolt/config.json", configDir),
		fmt.Sprintf("%s/.local/share/bolt/bolt.png", homeDir),
		fmt.Sprintf("%s/.local/share/bolt/error.log", homeDir),
		fmt.Sprintf("%s/.local/share/bolt/history.log", homeDir),
		fmt.Sprintf("%s/.local/share/applications/bolt.desktop", homeDir),
	}

	err = validateFiles(files, icon)
	if err != nil {
		return nil, fmt.Errorf("setup: couldn't validate default files:\n--> %w", err)
	}

	return files, nil
}

// validateFolders checks, if our folders in the user's config/cache dirs exists
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

// validateFiles checks, if config.json and search_cache.json exist
func validateFiles(filesToCheck []string, icon embed.FS) error {
	for _, file := range filesToCheck {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			continue
		}

		if strings.HasSuffix(file, "config.json") {
			resetConfig(file)
		} else if strings.HasSuffix(file, "bolt.desktop") {
			resetDotDesktop(file)
		} else if strings.HasSuffix(file, "bolt.png") {
			resetIcon(file, icon)
		}
	}

	return nil
}

// resetConfig resets the config file to the default settings
func resetConfig(configPath string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resetConfig: couldn't find the user's home dir:\n--> %w", err)
	}

	defaultConfig := Config{
		MaxCPUThreadPercentage:      0.25, // percentage of threads that may be used, always rounding the threads up
		ShortcutEnd:                 "space",
		DefaultDirsCacheUpdateTime:  120,  // in seconds
		ExtendedDirsCacheUpdateTime: 1800, // in seconds
		DefaultDirs: []string{
			fmt.Sprintf("%s/", homedir),
		},
		ExtendedDirs: []string{
			"/",
		},
		ExcludeFromDefaultDirs: Rules{
			Name: []string{},
			Path: []string{},
			Regex: []string{
				fmt.Sprintf(`^%s/\.[^/]+/?$`, homedir),
			},
		},
		ExcludeDirs: Rules{
			Name: []string{
				".git",
				"node_modules",
				"steamapps",
			},
			Path:  []string{},
			Regex: []string{},
		},
	}

	err = util.OverwriteJSON(configPath, true, defaultConfig)
	if err != nil {
		return fmt.Errorf("resetConfig: couldn't reset default config:\n--> %w", err)
	}

	return nil
}

// resetDotDesktop writes our default information into the provided .desktop
func resetDotDesktop(dotDesktopPath string) error {
	dotDesktopFile, err := os.OpenFile(dotDesktopPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("resetDotDesktop: couldn't open/create .desktop file:\n--> %w", err)
	}
	defer dotDesktopFile.Close()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resetDotDesktop: couldn't access the user's home dir:\n--> %w", err)
	}

	dotDesktopText := ("[Desktop Entry]\n" +
		"Name=Bolt\n" +
		"Type=Application\n" +
		"GenericName=Search Utility\n" +
		"Comment=Lightning fast search bar\n" +
		"Exec=env GDK_BACKEND=x11 /usr/local/bin/bolt\n" +
		fmt.Sprintf("Icon=%s/.local/share/bolt/bolt.png\n", homeDir) +
		"Terminal=false\n" +
		"Categories=Utility;")

	_, err = dotDesktopFile.WriteString(dotDesktopText)
	if err != nil {
		return fmt.Errorf("resetDotDesktop: couldn't write .desktop content to file:\n--> %w", err)
	}

	return nil
}

// resetIcon writes the build/appicon.png data into the provided .png
func resetIcon(iconPath string, icon embed.FS) error {
	iconData, err := icon.ReadFile("build/appicon.png")
	if err != nil {
		return fmt.Errorf("resetAppicon: read icon data:\n--> %w", err)
	}

	iconFile, err := os.OpenFile(iconPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("resetAppicon: couldn't open icon file:\n--> %w", err)
	}
	defer iconFile.Close()

	_, err = iconFile.Write(iconData)
	if err != nil {
		return fmt.Errorf("resetAppicon: couldn't write icon data to file:\n--> %w", err)
	}

	return nil
}
