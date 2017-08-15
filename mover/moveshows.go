package mover

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"tvrename/renamer"

	"github.com/nightlyone/lockfile"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

//MoveShowHandle stores the config and details for moving files
type MoveShowHandle struct {
	homeTvDirectory string
	shows           []string
}

//NewMoveShowHandler creates a new Handle to manage moving shows
func NewMoveShowHandler(homeTvDirectory string) *MoveShowHandle {
	shows := findDirectoryNames(homeTvDirectory)
	return &MoveShowHandle{homeTvDirectory: homeTvDirectory, shows: shows}
}

//MoveTvShowHome - takes the renamed file and moves it home
func (m *MoveShowHandle) MoveTvShowHome(t *renamer.TvShowDetails) {

	if t.Season == 0 {
		log.Printf("Not moving file as Season details were not found")
		return
	}

	showDirectory := m.findShowDirectory(t.Name)
	seasonDirectory := m.findSeasonDirectory(showDirectory, t.Season)

	lf, err := getFileLock(t.Path)
	if err != nil {
		return
	}
	os.Rename(t.Path, filepath.Join(seasonDirectory, t.ComputedName))

	lf.Unlock()
}

func getFileLock(path string) (lockfile.Lockfile, error) {
	lf, err := lockfile.New(path)
	if err != nil {
		fmt.Printf("Can not find file to lock")
		return lf, err
	}

	err = lf.TryLock()
	if err != nil {
		fmt.Printf("Cannot lock %q, reason: %v", lf, err)
		return lf, err
	}
	return lf, nil
}
func (m *MoveShowHandle) createDirectory(name string) string {
	os.Mkdir(m.homeTvDirectory+"/"+name, 07777)
	return name
}
func (m *MoveShowHandle) findShowDirectory(showName string) string {

	matches := fuzzy.RankFindFold(showName, m.shows)
	if len(matches) == 0 {
		showName = strings.Title(showName)
		m.createDirectory(showName)
		return showName
	}
	sort.Sort(matches)
	return matches[0].Target
}

func (m *MoveShowHandle) findSeasonDirectory(showDirectory string, season int) string {
	directoryNames := findDirectoryNames(m.homeTvDirectory + "/" + showDirectory)

	result := fuzzy.FindFold("Season "+strconv.Itoa(season), directoryNames)

	showPath := showDirectory + "/Season " + strconv.Itoa(season)
	if len(result) == 0 {
		m.createDirectory(showPath)
	}

	return m.homeTvDirectory + "/" + showPath
}

func findDirectoryNames(baseTvDirectory string) []string {
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
