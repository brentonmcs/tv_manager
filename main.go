package main

import (
	"flag"
	"fmt"
	"log"
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
	// Process events
	go func() {
		for {
			select {
			case event := <-w.Event:

				if !event.FileInfo.IsDir() {
					fmt.Printf("Not Directory")

					q.Push(simpleQueue.Item{Name: event.Name(),
						Fn: func() {
							details := renamer.GetTvShowDetails(event.Path)
							moveHandler.MoveTvShowHome(details)
						},
						StartAt: time.Now().Add(time.Second * 3)})
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

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
