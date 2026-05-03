package store

import (
	"os"
	"path/filepath"
)

func SettingsPath() string {
	return filepath.Join(os.Getenv("HOME"), ".claude", "settings.json")
}

func CachePath() string {
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	if xdgCacheHome == "" {
		xdgCacheHome = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return filepath.Join(xdgCacheHome, "proteus", "cache.json")
}

func BackupDir() string {
	return filepath.Join(os.Getenv("HOME"), ".claude", "proteus-backups")
}

func ProteusConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(xdgConfigHome, "proteus")
}

func LaunchProfileConfigDir(profile string) string {
	return filepath.Join(ProteusConfigDir(), "claude", "profiles", profile)
}

func LaunchProfileSettingsPath(profile string) string {
	return filepath.Join(LaunchProfileConfigDir(profile), "settings.json")
}
