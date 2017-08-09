package main

import (
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

	moveHandler = mover.NewMoveShowHandler("/Users/brentonmcsweyn/TvShows")
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

	if err := w.Add("/Users/brentonmcsweyn/Downloads"); err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting Queue")
	q.Run()

	log.Println("Starting Watcher")
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
