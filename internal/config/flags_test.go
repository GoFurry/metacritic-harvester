package config

import "testing"

func TestBuildListCommandConfig(t *testing.T) {
	t.Parallel()

	cfg, err := BuildListCommandConfig(ListCommandOptions{
		Category:    "movie",
		Metric:      "userscore",
		Year:        "2011:2014",
		Network:     "netflix,max",
		Genre:       "drama,thriller",
		ReleaseType: "coming-soon,in-theaters",
		Pages:       3,
		DBPath:      "output/test.db",
		Debug:       true,
		MaxRetries:  2,
		Proxies:     "http://127.0.0.1:7897",
	})
	if err != nil {
		t.Fatalf("BuildListCommandConfig() error = %v", err)
	}

	if cfg.Task.Category != "movie" || cfg.Task.Metric != "userscore" {
		t.Fatalf("unexpected task: %+v", cfg.Task)
	}
	if cfg.Task.MaxPages != 3 || !cfg.Debug || cfg.MaxRetries != 2 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.Task.Filter.ReleaseYearMin == nil || *cfg.Task.Filter.ReleaseYearMin != 2011 {
		t.Fatalf("expected release year min 2011, got %+v", cfg.Task.Filter.ReleaseYearMin)
	}
	if cfg.Task.Filter.ReleaseYearMax == nil || *cfg.Task.Filter.ReleaseYearMax != 2014 {
		t.Fatalf("expected release year max 2014, got %+v", cfg.Task.Filter.ReleaseYearMax)
	}
	if len(cfg.Task.Filter.Networks) != 2 || len(cfg.Task.Filter.Genres) != 2 || len(cfg.Task.Filter.ReleaseTypes) != 2 {
		t.Fatalf("unexpected filter parsing: %+v", cfg.Task.Filter)
	}
	if len(cfg.ProxyURLs) != 1 {
		t.Fatalf("expected one proxy, got %d", len(cfg.ProxyURLs))
	}
}

func TestBuildListCommandConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []ListCommandOptions{
		{Category: "bad", Metric: "metascore", Pages: 1, DBPath: "output/test.db"},
		{Category: "game", Metric: "bad", Pages: 1, DBPath: "output/test.db"},
		{Category: "game", Metric: "metascore", Pages: 0, DBPath: "output/test.db"},
		{Category: "game", Metric: "metascore", Pages: 1, DBPath: " ", Proxies: "http://127.0.0.1:7897"},
		{Category: "game", Metric: "metascore", Pages: 1, DBPath: "output/test.db", Proxies: "bad-proxy"},
		{Category: "game", Metric: "metascore", Pages: 1, DBPath: "output/test.db", Year: "2014:2011"},
		{Category: "game", Metric: "metascore", Pages: 1, DBPath: "output/test.db", Year: "2011-2014"},
		{Category: "movie", Metric: "metascore", Pages: 1, DBPath: "output/test.db", Platform: "pc"},
		{Category: "tv", Metric: "metascore", Pages: 1, DBPath: "output/test.db", ReleaseType: "coming-soon"},
		{Category: "game", Metric: "metascore", Pages: 1, DBPath: "output/test.db", Network: "netflix"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Category+"-"+tt.Metric, func(t *testing.T) {
			t.Parallel()
			if _, err := BuildListCommandConfig(tt); err == nil {
				t.Fatalf("expected error for %+v", tt)
			}
		})
	}
}
