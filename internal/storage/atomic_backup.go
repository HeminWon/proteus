package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/HeminWon/proteus/internal/cli"
)

func writeFileAtomic(filePath string, content []byte) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmpPath := fmt.Sprintf("%s.tmp-%d-%d", filePath, os.Getpid(), time.Now().UnixMilli())
	if err := os.WriteFile(tmpPath, content, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, filePath)
}

func formatTimestamp(t time.Time) string {
	return t.Format("20060102_150405")
}

func CreateBackupIfNeeded(settingsExists bool) (string, error) {
	if !settingsExists {
		return "", nil
	}

	if err := os.MkdirAll(BackupDir(), 0o755); err != nil {
		return "", err
	}

	backupPath := filepath.Join(BackupDir(), fmt.Sprintf("settings-%s.json", formatTimestamp(time.Now())))
	input, err := os.ReadFile(SettingsPath())
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(backupPath, input, 0o644); err != nil {
		return "", err
	}

	if err := CleanupOldBackups(); err != nil {
		return "", err
	}
	return backupPath, nil
}

func CleanupOldBackups() error {
	entries, err := os.ReadDir(BackupDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	type backupInfo struct {
		path    string
		mtimeMs int64
	}
	backups := make([]backupInfo, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) != len("settings-20060102_150405.json") || name[:9] != "settings-" || name[len(name)-5:] != ".json" {
			continue
		}
		full := filepath.Join(BackupDir(), name)
		info, statErr := os.Stat(full)
		if statErr != nil {
			continue
		}
		backups = append(backups, backupInfo{path: full, mtimeMs: info.ModTime().UnixMilli()})
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].mtimeMs > backups[j].mtimeMs
	})

	if len(backups) <= core.MaxBackups {
		return nil
	}
	for _, stale := range backups[core.MaxBackups:] {
		if err := os.Remove(stale.path); err != nil {
			return err
		}
	}
	return nil
}

func RestoreFromBackup(backupPath string) error {
	if backupPath == "" {
		return nil
	}
	input, err := os.ReadFile(backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.WriteFile(SettingsPath(), input, 0o644)
}
