package main

import (
	"testing"
)

/**
// !17. track17.mp3
// 07. Di-Rect - I Just Can't Stand.mp3
// 07.Di-Rect - I Just Can't Stand.mp3
// 10_Forever Gone.mp3
// 10 Forever Gone.mp3
// 11-Lady.mp3
// 11 - Lady.mp3

// Track  4.mp3
// Track  5.mp3

// Portal2-16-Hard_Sunshine.mp3
// Portal2-17-I_Am_Different.mp3
*/
func TestAnalyzeListEntry(t *testing.T) {

	var tests = []struct {
		s, want string
		good    bool
	}{

		{"07. Di-Rect - I Just Can't Stand.mp3", "Di-Rect - I Just Can't Stand", true},
		{" 07.Di-Rect - I Just Can't Stand.mp3", "Di-Rect - I Just Can't Stand", true},
		{"10_Forever Gone.mp3", "Forever Gone", true},
		{"10 Forever Gone.mp3", "Forever Gone", true},
		{"11-Lady.mp3", "Lady", true},
		{"11 - Lady.mp3", "Lady", true},
		{"11 - Smith's son", "Smith's son", true},
		{"Smith's son", "Smith's son", true},
		{"Track 4.mp3", "Track", true},
		{"Track  010.mp3", "Track", true},
		{"Track-10", "Track", true},
		{"Photo_5", "Photo", true},
		{"Photo_16.jpg", "Photo", true},
		{".git", "", false},
	}

	for _, c := range tests {

		good, got := analyzeListEntry(c.s)

		if good != c.good {
			t.Errorf("While analyzing %q got error flag %v, expected %v", c.s, good, c.good)
		} else if got.stem != c.want {
			t.Log(got.name)
			t.Errorf("%q(%q) ==> (%q), want %q", c.s, got.name, got.stem, c.want)
		}

	}

}
