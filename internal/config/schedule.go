package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ScheduleFile struct {
	Timezone string        `yaml:"timezone"`
	Jobs     []ScheduleJob `yaml:"jobs"`
}

type ScheduleJob struct {
	Name        string `yaml:"name"`
	Cron        string `yaml:"cron"`
	BatchFile   string `yaml:"batch_file"`
	Enabled     *bool  `yaml:"enabled"`
	Concurrency *int   `yaml:"concurrency"`
}

func LoadScheduleFile(path string) (ScheduleFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ScheduleFile{}, fmt.Errorf("read schedule file: %w", err)
	}

	var scheduleFile ScheduleFile
	if err := yaml.Unmarshal(content, &scheduleFile); err != nil {
		return ScheduleFile{}, fmt.Errorf("parse schedule file: %w", err)
	}
	if len(scheduleFile.Jobs) == 0 {
		return ScheduleFile{}, fmt.Errorf("schedule file must include at least one job")
	}

	baseDir := filepath.Dir(path)
	for idx := range scheduleFile.Jobs {
		job := &scheduleFile.Jobs[idx]
		job.Name = strings.TrimSpace(job.Name)
		job.Cron = strings.TrimSpace(job.Cron)
		job.BatchFile = strings.TrimSpace(job.BatchFile)

		if job.Name == "" {
			return ScheduleFile{}, fmt.Errorf("job %d: name is required", idx+1)
		}
		if job.Cron == "" {
			return ScheduleFile{}, fmt.Errorf("job %s: cron is required", job.Name)
		}
		if job.BatchFile == "" {
			return ScheduleFile{}, fmt.Errorf("job %s: batch_file is required", job.Name)
		}
		if !filepath.IsAbs(job.BatchFile) {
			job.BatchFile = filepath.Clean(filepath.Join(baseDir, job.BatchFile))
		}
	}

	return scheduleFile, nil
}

func (j ScheduleJob) IsEnabled() bool {
	if j.Enabled == nil {
		return true
	}
	return *j.Enabled
}
