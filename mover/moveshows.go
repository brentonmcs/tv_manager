package mover

import (
	"io/ioutil"
	"log"
	"sort"
	"tvrename/renamer"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

//MoveShowHandle stores the config and details for moving files
type MoveShowHandle struct {
	watchDirectory string
	shows          []string
}

//NewMoveShowHandler creates a new Handle to manage moving shows
func NewMoveShowHandler(watchDirectory string) *MoveShowHandle {
	shows := findShows(watchDirectory)
	return &MoveShowHandle{watchDirectory: watchDirectory, shows: shows}
}

//MoveTvShowHome - takes the renamed file and moves it home
func (m *MoveShowHandle) MoveTvShowHome(baseTvDirectory string, t *renamer.TvShowDetails) {

}

func (m *MoveShowHandle) findShowDirectory(showName string) string {

	matches := fuzzy.RankFindFold(showName, m.shows)
	if len(matches) == 0 {
		return ""
	}
	sort.Sort(matches)
	return matches[0].Target
}

func findShows(baseTvDirectory string) []string {
	files, err := ioutil.ReadDir(baseTvDirectory)
	if err != nil {
		log.Fatal(err)
	}
	var result []string
	for _, f := range files {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}
	return result
}
