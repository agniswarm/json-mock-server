package notifier

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// watchFileChanges watches for changes in the JSON file
func WatchFileChanges(filePath string, reload chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(filepath.Dir(filePath))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 && event.Name == filePath {
				log.Println("File changed, reloading routes...")
				reload <- true
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error watching file:", err)
		}
	}
}
