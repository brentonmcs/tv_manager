package main

import (
	"log"
	"tvrename/mover"
	"tvrename/renamer"

	"github.com/howeyc/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					log.Println("event:", ev)
					tvShowDetails := renamer.GetTvShowDetails(ev.Name).RenameFile()
					mover.MoveTvShowHome(tvShowDetails)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("/Users/brentonmcsweyn/tools")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}
