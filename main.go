package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"tvrename/mover"
	"tvrename/renamer"
	"tvrename/simpleQueue"

	"github.com/radovskyb/watcher"
)

var moveHandler *mover.MoveShowHandle

func main() {

	w := watcher.New()
	w.IgnoreHiddenFiles(true)
	w.FilterOps(watcher.Write)
	q := simpleQueue.NewQueue(time.Millisecond * 400)

	watchDir := flag.String("watchDir", "", "the folder to watch for incoming files")
	homeTvDir := flag.String("homeTvDir", "", "the root folder for all TV Shows")
	processOnly := flag.Bool("processOnly", false, "process only files in the folder - do not watch")

	flag.Parse()
	if *watchDir == "" {
		log.Fatal("Missing watchDir  command flags")
		return
	}

	if *homeTvDir == "" {
		log.Fatal("Missing homeTvDir command flags")
		return
	}
	moveHandler = mover.NewMoveShowHandler(*homeTvDir)

	processExistingItems(*watchDir)
	// Process events
	go func() {
		for {
			select {
			case event := <-w.Event:

				if !event.FileInfo.IsDir() {
					fmt.Printf("Not Directory")

					details := renamer.GetTvShowDetails(event.Path)
					moveHandler.MoveTvShowHome(details)
					// q.Push(simpleQueue.Item{Name: event.Name(),
					// 	Fn: func() {

					// 	},
					// 	StartAt: time.Now().Add(time.Second * 3)})
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if !*processOnly {
		if err := w.Add(*watchDir); err != nil {
			log.Fatalln(err)
		}

		log.Println("Starting Queue")
		q.Run()

		log.Println("Starting Watcher")
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}
}

func processExistingItems(watchDir string) {
	files, err := ioutil.ReadDir(watchDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		if !strings.HasPrefix(f.Name(), ".") {

			path := watchDir + "/" + f.Name()
			if f.IsDir() {
				processExistingItems(path)
			} else {
				details := renamer.GetTvShowDetails(path)
				moveHandler.MoveTvShowHome(details)
			}
		}
	}
}
