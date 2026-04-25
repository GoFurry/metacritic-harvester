package domain

import "testing"

func TestFilterKeyIsStable(t *testing.T) {
	t.Parallel()

	minYear := 2011
	maxYear := 2014

	a := Filter{
		ReleaseYearMin: &minYear,
		ReleaseYearMax: &maxYear,
		Platforms:      []string{"ps5", "pc", "pc"},
		Networks:       []string{"max", "netflix"},
		Genres:         []string{"rpg", "action"},
		ReleaseTypes:   []string{"in-theaters", "coming-soon"},
	}
	b := Filter{
		ReleaseYearMin: &minYear,
		ReleaseYearMax: &maxYear,
		Platforms:      []string{"pc", "ps5"},
		Networks:       []string{"netflix", "max"},
		Genres:         []string{"action", "rpg"},
		ReleaseTypes:   []string{"coming-soon", "in-theaters"},
	}

	if a.Key() != b.Key() {
		t.Fatalf("expected identical filter keys, got %q and %q", a.Key(), b.Key())
	}
}
