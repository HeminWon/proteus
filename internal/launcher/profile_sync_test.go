package launcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanProfileConfigSyncLinksWhitelistEntries(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")

	for _, name := range sharedConfigWhitelist {
		source := filepath.Join(home, ".claude", name)
		if err := os.MkdirAll(source, 0o755); err != nil {
			t.Fatalf("mkdir source %s: %v", name, err)
		}
	}

	entries, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}

	for _, name := range sharedConfigWhitelist {
		entry := findSyncEntry(t, entries, name)
		if entry.Status != SyncStatusLinked {
			t.Fatalf("entry %s status = %s, want %s", name, entry.Status, SyncStatusLinked)
		}
	}
	claude := findSyncEntry(t, entries, "CLAUDE.md")
	if claude.Status != SyncStatusDisabled {
		t.Fatalf("CLAUDE.md status = %s, want %s", claude.Status, SyncStatusDisabled)
	}
}

func TestPlanProfileConfigSyncReusedWhenSymlinkMatches(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		t.Fatalf("mkdir profile dir: %v", err)
	}

	source := filepath.Join(home, ".claude", "commands")
	target := filepath.Join(profileDir, "commands")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.Symlink(source, target); err != nil {
		t.Fatalf("symlink target: %v", err)
	}

	entries, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}

	entry := findSyncEntry(t, entries, "commands")
	if entry.Status != SyncStatusReused {
		t.Fatalf("commands status = %s, want %s", entry.Status, SyncStatusReused)
	}
}

func TestPlanProfileConfigSyncSkipsMissingAndKeepsRuntimeOut(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")

	if err := os.MkdirAll(filepath.Join(home, ".claude", "tasks"), 0o755); err != nil {
		t.Fatalf("mkdir runtime tasks dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(home, ".claude", "history.jsonl"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write runtime history file: %v", err)
	}

	entries, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}

	for _, name := range sharedConfigWhitelist {
		entry := findSyncEntry(t, entries, name)
		if entry.Status != SyncStatusSkippedMissing {
			t.Fatalf("entry %s status = %s, want %s", name, entry.Status, SyncStatusSkippedMissing)
		}
		if strings.Contains(entry.TargetPath, "tasks") || strings.Contains(entry.TargetPath, "history.jsonl") {
			t.Fatalf("runtime entry leaked into sync plan: %+v", entry)
		}
	}
}

func TestPlanProfileConfigSyncConflictsOnUnexpectedDestinationKinds(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		t.Fatalf("mkdir profile dir: %v", err)
	}

	tests := []struct {
		name  string
		setup func(path string) error
	}{
		{
			name: "regular-file",
			setup: func(path string) error {
				return os.WriteFile(path, []byte("x"), 0o644)
			},
		},
		{
			name: "directory",
			setup: func(path string) error {
				return os.MkdirAll(path, 0o755)
			},
		},
		{
			name: "wrong-symlink",
			setup: func(path string) error {
				wrong := filepath.Join(home, ".claude", "wrong-skills")
				if err := os.MkdirAll(wrong, 0o755); err != nil {
					return err
				}
				return os.Symlink(wrong, path)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.RemoveAll(profileDir); err != nil {
				t.Fatalf("reset profile dir: %v", err)
			}
			if err := os.MkdirAll(profileDir, 0o755); err != nil {
				t.Fatalf("recreate profile dir: %v", err)
			}

			source := filepath.Join(home, ".claude", "skills")
			if err := os.MkdirAll(source, 0o755); err != nil {
				t.Fatalf("mkdir source: %v", err)
			}
			conflict := filepath.Join(profileDir, "skills")
			if err := tt.setup(conflict); err != nil {
				t.Fatalf("setup conflict: %v", err)
			}

			entries, err := PlanProfileConfigSync(profileDir, false)
			if err != nil {
				t.Fatalf("PlanProfileConfigSync error = %v", err)
			}
			if st := findSyncEntry(t, entries, "skills").Status; st != SyncStatusConflict {
				t.Fatalf("skills status = %s, want %s", st, SyncStatusConflict)
			}

			err = ApplyProfileConfigSync(entries, profileDir)
			if err == nil {
				t.Fatalf("expected apply conflict error")
			}
			if got := err.Error(); !strings.Contains(got, conflict) {
				t.Fatalf("error = %q, want contains %q", got, conflict)
			}
		})
	}
}

func TestPlanProfileConfigSyncClaudeMdToggle(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")

	claudeSource := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.WriteFile(claudeSource, []byte("rules"), 0o644); err != nil {
		t.Fatalf("write source CLAUDE.md: %v", err)
	}

	entriesOn, err := PlanProfileConfigSync(profileDir, true)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync enabled error = %v", err)
	}
	if st := findSyncEntry(t, entriesOn, "CLAUDE.md").Status; st != SyncStatusLinked {
		t.Fatalf("enabled CLAUDE.md status = %s, want %s", st, SyncStatusLinked)
	}

	entriesOff, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync disabled error = %v", err)
	}
	if st := findSyncEntry(t, entriesOff, "CLAUDE.md").Status; st != SyncStatusDisabled {
		t.Fatalf("disabled CLAUDE.md status = %s, want %s", st, SyncStatusDisabled)
	}
}

func TestPlanProfileConfigSyncClaudeMdDisabledExistingRetained(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		t.Fatalf("mkdir profile dir: %v", err)
	}

	claudeSource := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.WriteFile(claudeSource, []byte("rules"), 0o644); err != nil {
		t.Fatalf("write source CLAUDE.md: %v", err)
	}
	claudeTarget := filepath.Join(profileDir, "CLAUDE.md")
	if err := os.Symlink(claudeSource, claudeTarget); err != nil {
		t.Fatalf("symlink target CLAUDE.md: %v", err)
	}

	entries, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}
	if st := findSyncEntry(t, entries, "CLAUDE.md").Status; st != SyncStatusDisabledExisting {
		t.Fatalf("disabled existing CLAUDE.md status = %s, want %s", st, SyncStatusDisabledExisting)
	}
}

func TestApplyProfileConfigSyncCreatesExpectedSymlinks(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	for _, name := range sharedConfigWhitelist {
		source := filepath.Join(home, ".claude", name)
		if err := os.MkdirAll(source, 0o755); err != nil {
			t.Fatalf("mkdir source %s: %v", name, err)
		}
	}
	claudeSource := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.WriteFile(claudeSource, []byte("rules"), 0o644); err != nil {
		t.Fatalf("write source CLAUDE.md: %v", err)
	}

	entries, err := PlanProfileConfigSync(profileDir, true)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}
	if err := ApplyProfileConfigSync(entries, profileDir); err != nil {
		t.Fatalf("ApplyProfileConfigSync error = %v", err)
	}

	for _, name := range append(append([]string{}, sharedConfigWhitelist...), "CLAUDE.md") {
		target := filepath.Join(profileDir, name)
		info, err := os.Lstat(target)
		if err != nil {
			t.Fatalf("lstat %s: %v", target, err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("%s is not a symlink", target)
		}
	}
}

func TestPlanAndApplyStayIdempotent(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	for _, name := range sharedConfigWhitelist {
		source := filepath.Join(home, ".claude", name)
		if err := os.MkdirAll(source, 0o755); err != nil {
			t.Fatalf("mkdir source %s: %v", name, err)
		}
	}

	entries, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync first error = %v", err)
	}
	if err := ApplyProfileConfigSync(entries, profileDir); err != nil {
		t.Fatalf("ApplyProfileConfigSync first error = %v", err)
	}

	replanned, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync second error = %v", err)
	}
	for _, name := range sharedConfigWhitelist {
		if st := findSyncEntry(t, replanned, name).Status; st != SyncStatusReused {
			t.Fatalf("entry %s status after apply = %s, want %s", name, st, SyncStatusReused)
		}
	}
}

func TestPlanProfileConfigSyncDryRunNoSideEffects(t *testing.T) {
	home := setupProfileSyncEnv(t)
	profileDir := filepath.Join(home, ".config", "proteus", "claude", "profiles", "default")
	if err := os.MkdirAll(filepath.Join(home, ".claude", "commands"), 0o755); err != nil {
		t.Fatalf("mkdir source commands: %v", err)
	}

	if _, err := os.Stat(profileDir); !os.IsNotExist(err) {
		t.Fatalf("expected profile dir not to exist before planning")
	}

	_, err := PlanProfileConfigSync(profileDir, false)
	if err != nil {
		t.Fatalf("PlanProfileConfigSync error = %v", err)
	}

	if _, err := os.Stat(profileDir); !os.IsNotExist(err) {
		t.Fatalf("expected profile dir to remain absent after planning")
	}
}

func setupProfileSyncEnv(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0o755); err != nil {
		t.Fatalf("mkdir ~/.claude: %v", err)
	}
	return home
}

func findSyncEntry(t *testing.T, entries []ProfileSyncEntry, name string) ProfileSyncEntry {
	t.Helper()
	for _, entry := range entries {
		if entry.Name == name {
			return entry
		}
	}
	t.Fatalf("sync entry %s not found", name)
	return ProfileSyncEntry{}
}
