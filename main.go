package main

import (
	"flag"
	"log"
	"log/syslog"
	"tvrename/mover"
	"tvrename/renamer"
)

var moveHandler *mover.MoveShowHandle

func main() {

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "tvrename")
	if e == nil {
		log.SetOutput(logwriter)
	}

	showPath := flag.String("showPath", "", "the movie file path")
	homeTvDir := flag.String("homeTvDir", "", "the root folder for all TV Shows")
	flag.Parse()

	if *showPath == "" {
		log.Fatal("Missing showPath command flag")
		return
	}

	if *homeTvDir == "" {
		log.Fatal("Missing homeTvDir command flags")
		return
	}
	moveHandler = mover.NewMoveShowHandler(*homeTvDir)

	log.Printf("Processing filename : %s", *showPath)
	moveHandler.MoveTvShowHome(renamer.GetTvShowDetails(*showPath))
}
