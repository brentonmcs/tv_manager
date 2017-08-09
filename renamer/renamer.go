package renamer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//TvShowDetails holds the information about the show
type TvShowDetails struct {
	name         string
	path         string
	season       int
	episode      int
	extension    string
	computedName string
}

type fileDetails struct {
	filename  string
	path      string
	extension string
}

type parsedFileDetails struct {
	season  int
	episode int
	name    string
}

func (p *parsedFileDetails) validSeasonFound() bool {
	return p.season == 0 && p.episode == 0
}

//NewTvShowDetails constuctor
func NewTvShowDetails(parsedFileDetails parsedFileDetails, extension, path string) *TvShowDetails {
	computedName := fmt.Sprintf("%v S%02dE%02d.%v", strings.Title(parsedFileDetails.name), parsedFileDetails.season, parsedFileDetails.episode, extension)
	return &TvShowDetails{name: parsedFileDetails.name, season: parsedFileDetails.season, episode: parsedFileDetails.episode, extension: extension, computedName: computedName, path: path}
}

//NewTvShowDetailsWithComputedPath constuctor - updates path to computed path
func NewTvShowDetailsWithComputedPath(t *TvShowDetails) *TvShowDetails {

	fName := filepath.Dir(t.path)
	newPath := fName + "/" + t.computedName
	return &TvShowDetails{name: t.name, path: newPath, computedName: t.computedName, extension: t.extension, season: t.season, episode: t.episode}
}

//RenameFile takes the show details and renames them
func (t *TvShowDetails) RenameFile() *TvShowDetails {
	err := os.Rename(t.path, filepath.Dir(t.path)+"/"+t.computedName)

	if err != nil {
		log.Fatal("Failed to rename")
	}

	return NewTvShowDetailsWithComputedPath(t)
}

//GetTvShowDetails renames the file so it's a clean TV Show Name
func GetTvShowDetails(path string) *TvShowDetails {
	return findDetails(getFileDetails(path))
}

func getFileDetails(path string) fileDetails {

	fName := filepath.Base(path)
	extName := filepath.Ext(path)
	bName := fName[:len(fName)-len(extName)]

	return fileDetails{filename: bName, path: path, extension: extName}
}

func findName(filename string) string {

	if strings.Contains(filename, ".") {
		return strings.Split(filename, ".")[0]
	}
	return filename
}

func findDetails(file fileDetails) *TvShowDetails {

	parsedFileDetails := extractViaRegex("S\\d+E\\d+", extraDigitsSE, file.filename)

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d+X\\d+", extraDigitsSE, file.filename)
	}

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d\\d\\d\\d", extraDigits4d, file.filename)
	}

	if parsedFileDetails.validSeasonFound() {
		parsedFileDetails = extractViaRegex("\\d\\d\\d", extractThreeDigit, file.filename)
	}
	return NewTvShowDetails(parsedFileDetails, file.extension, file.path)
}

func removeDot(filename string) string {
	return strings.Trim(strings.Replace(filename, ".", " ", -1), " ")
}
