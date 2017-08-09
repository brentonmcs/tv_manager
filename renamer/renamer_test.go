package renamer

import "testing"

type TestName struct {
	path     string
	expected string
}

type TestDetails struct {
	path    string
	season  int
	episode int
}

func TestRenameName(t *testing.T) {
	testCase := []TestName{
		TestName{path: "/User/Home/test.local", expected: "test"},
		TestName{path: "/User/Home/test.S01E01.HDTV.local", expected: "test"},
		TestName{path: "/User/Home/the.big.bang.theory.s9e13.hdtv-lol.mp4", expected: "the big bang theory"},
	}

	for _, c := range testCase {
		result := GetTvShowDetails(c.path)

		t.Logf("Rename Info - expected %v, result %v", c.expected, result.name)
		if c.expected != result.name {
			t.Fatalf("Rename is not correct - expected %v, result %v", c.expected, result.name)
		}
	}

}

func TestRenameSeasonDetails(t *testing.T) {
	testCase := []TestDetails{
		TestDetails{path: "/User/Home/test.S01E01.HDTV.local", season: 1, episode: 1},
		TestDetails{path: "/User/Home/test.s01e01.HDTV.local", season: 1, episode: 1},
		TestDetails{path: "/User/Home/test.02X11.HDTV.local", season: 2, episode: 11},
		TestDetails{path: "/User/Home/test.02x11.HDTV.local", season: 2, episode: 11},
		TestDetails{path: "/User/Home/test.0212.HDTV.local", season: 2, episode: 12},
		TestDetails{path: "/User/Home/test.212.HDTV.local", season: 2, episode: 12},
		TestDetails{path: "/User/Home/the.big.bang.theory.s9e13.hdtv-lol.mp4", season: 9, episode: 13},
	}

	for _, c := range testCase {
		result := GetTvShowDetails(c.path)

		if c.season != 0 && c.season != result.season {
			t.Fatalf("Rename Season is not correct - expected %v, result %v", c.season, result.season)
		}

		if c.episode != 0 && c.episode != result.episode {
			t.Fatalf("Rename episode is not correct - expected %v, result %v", c.season, result.episode)
		}
	}

}

func TestComputedName(t *testing.T) {
	result := NewTvShowDetails(parsedFileDetails{episode: 1, season: 1, name: "lower"}, "mp3", "")
	expected := "Lower S01E01.mp3"

	t.Logf("%v", result.computedName)
	if result.computedName != expected {
		t.Fatalf("Computed Name is not correct - expected %v, result %v", result.computedName, expected)
	}
	result = NewTvShowDetails(parsedFileDetails{episode: 1, season: 1, name: "the lower name"}, "mp3", "")
	expected = "The Lower Name S01E01.mp3"

	t.Logf("%v", result.computedName)
	if result.computedName != expected {
		t.Fatalf("Computed Name is not correct - expected %v, result %v", result.computedName, expected)
	}
}

func TestNewDetailsWithComputedPath(t *testing.T) {
	details := NewTvShowDetails(parsedFileDetails{episode: 1, season: 1, name: "lower"}, "mp3", "/home/brenton/lower.mp3")

	result := NewTvShowDetailsWithComputedPath(details)

	expected := "/home/brenton/Lower S01E01.mp3"

	if result.path != expected {
		t.Fatalf("Path Name is not correct - expected %v, result %v", result.path, expected)
	}
}
