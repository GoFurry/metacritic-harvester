package cli

import (
	"context"
	"testing"
)

func TestCrawlScheduleCommandPassesFile(t *testing.T) {
	t.Parallel()

	called := false
	cmd := newCrawlScheduleCommandWithRunner(func(_ context.Context, filePath string) error {
		called = true
		if filePath != "schedule.yaml" {
			t.Fatalf("unexpected file path %q", filePath)
		}
		return nil
	})

	cmd.SetArgs([]string{"--file=schedule.yaml"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !called {
		t.Fatal("expected runner to be called")
	}
}
