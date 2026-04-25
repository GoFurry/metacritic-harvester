package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadScheduleFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schedulePath := filepath.Join(dir, "schedule.yaml")
	if err := os.WriteFile(schedulePath, []byte(`
timezone: Asia/Shanghai
jobs:
  - name: nightly
    cron: "0 0 * * *"
    batch_file: ./tasks.yaml
`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	file, err := LoadScheduleFile(schedulePath)
	if err != nil {
		t.Fatalf("LoadScheduleFile() error = %v", err)
	}
	if file.Timezone != "Asia/Shanghai" {
		t.Fatalf("unexpected timezone %q", file.Timezone)
	}
	if len(file.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(file.Jobs))
	}
	if !filepath.IsAbs(file.Jobs[0].BatchFile) {
		t.Fatalf("expected batch file to resolve to absolute path, got %q", file.Jobs[0].BatchFile)
	}
	if !file.Jobs[0].IsEnabled() {
		t.Fatal("expected job to be enabled by default")
	}
}

func TestLoadScheduleFileRejectsEmptyJobs(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schedulePath := filepath.Join(dir, "schedule.yaml")
	if err := os.WriteFile(schedulePath, []byte("timezone: UTC\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := LoadScheduleFile(schedulePath); err == nil {
		t.Fatal("expected error for empty jobs")
	}
}
