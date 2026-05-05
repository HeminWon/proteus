package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ProfileSyncStatus string

const (
	SyncStatusLinked           ProfileSyncStatus = "linked"
	SyncStatusReused           ProfileSyncStatus = "reused"
	SyncStatusSkippedMissing   ProfileSyncStatus = "skipped-missing"
	SyncStatusConflict         ProfileSyncStatus = "conflict"
	SyncStatusDisabled         ProfileSyncStatus = "disabled"
	SyncStatusDisabledExisting ProfileSyncStatus = "disabled-existing"
)

type ProfileSyncEntry struct {
	Name       string
	SourcePath string
	TargetPath string
	Status     ProfileSyncStatus
}

var sharedConfigWhitelist = []string{"commands", "skills", "plugins", "agents", "ide"}

func globalClaudeDir() string {
	return filepath.Join(os.Getenv("HOME"), ".claude")
}

func PlanProfileConfigSync(profileConfigDir string, shareClaudeMD bool) ([]ProfileSyncEntry, error) {
	entries := make([]ProfileSyncEntry, 0, len(sharedConfigWhitelist)+1)
	sourceRoot := globalClaudeDir()

	for _, name := range sharedConfigWhitelist {
		sourcePath := filepath.Join(sourceRoot, name)
		targetPath := filepath.Join(profileConfigDir, name)
		status, err := plannedSyncStatus(sourcePath, targetPath)
		if err != nil {
			return nil, err
		}
		entries = append(entries, ProfileSyncEntry{
			Name:       name,
			SourcePath: sourcePath,
			TargetPath: targetPath,
			Status:     status,
		})
	}

	claudeSource := filepath.Join(sourceRoot, "CLAUDE.md")
	claudeTarget := filepath.Join(profileConfigDir, "CLAUDE.md")
	if !shareClaudeMD {
		status := SyncStatusDisabled
		if exists, err := pathExists(claudeTarget); err != nil {
			return nil, err
		} else if exists {
			status = SyncStatusDisabledExisting
		}
		entries = append(entries, ProfileSyncEntry{
			Name:       "CLAUDE.md",
			SourcePath: claudeSource,
			TargetPath: claudeTarget,
			Status:     status,
		})
		return entries, nil
	}

	status, err := plannedSyncStatus(claudeSource, claudeTarget)
	if err != nil {
		return nil, err
	}
	entries = append(entries, ProfileSyncEntry{
		Name:       "CLAUDE.md",
		SourcePath: claudeSource,
		TargetPath: claudeTarget,
		Status:     status,
	})

	return entries, nil
}

func plannedSyncStatus(sourcePath, targetPath string) (ProfileSyncStatus, error) {
	if exists, err := pathExists(sourcePath); err != nil {
		return "", err
	} else if !exists {
		return SyncStatusSkippedMissing, nil
	}

	info, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return SyncStatusLinked, nil
		}
		return "", err
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return SyncStatusConflict, nil
	}

	matches, err := symlinkPointsTo(targetPath, sourcePath)
	if err != nil {
		return "", err
	}
	if matches {
		return SyncStatusReused, nil
	}

	return SyncStatusConflict, nil
}

func ApplyProfileConfigSync(entries []ProfileSyncEntry, profileConfigDir string) error {
	if err := os.MkdirAll(profileConfigDir, 0o755); err != nil {
		return err
	}

	for _, entry := range entries {
		switch entry.Status {
		case SyncStatusConflict, SyncStatusLinked:
			if err := ensureExpectedSymlink(entry.SourcePath, entry.TargetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func ensureExpectedSymlink(sourcePath, targetPath string) error {
	if exists, err := pathExists(sourcePath); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("profile sync source missing at %s", sourcePath)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}

	info, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.Symlink(sourcePath, targetPath)
		}
		return err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		matches, matchErr := symlinkPointsTo(targetPath, sourcePath)
		if matchErr != nil {
			return matchErr
		}
		if matches {
			return nil
		}
	}

	if err := backupExistingPath(targetPath); err != nil {
		return err
	}
	return os.Symlink(sourcePath, targetPath)
}

func backupExistingPath(path string) error {
	stamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup.%s", path, stamp)
	for i := 1; ; i++ {
		if exists, err := pathExists(backupPath); err != nil {
			return err
		} else if !exists {
			return os.Rename(path, backupPath)
		}
		backupPath = fmt.Sprintf("%s.backup.%s-%d", path, stamp, i)
	}
}

func symlinkPointsTo(linkPath, expectedTarget string) (bool, error) {
	actualTarget, err := os.Readlink(linkPath)
	if err != nil {
		return false, err
	}

	actualAbs := actualTarget
	if !filepath.IsAbs(actualAbs) {
		actualAbs = filepath.Join(filepath.Dir(linkPath), actualTarget)
	}

	actualAbs = filepath.Clean(actualAbs)
	expectedAbs := filepath.Clean(expectedTarget)
	return actualAbs == expectedAbs, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
