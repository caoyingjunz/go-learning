package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// more info https://github.com/fsnotify/fsnotify
// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/pluginmanager/pluginwatcher/plugin_watcher.go#L60

func handleCreateEvent(event fsnotify.Event) error {
	log.Println("Add file", event.Name)
	return nil
}

func handleDeleteEvent(event fsnotify.Event) error {
	log.Println("Delete file", event.Name)
	return nil
}

func main() {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer fsWatcher.Close()

	go func(fsWatcher *fsnotify.Watcher) {
		for {
			select {
			case event := <-fsWatcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					_ = handleCreateEvent(event)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					handleDeleteEvent(event)
				}
			case err = <-fsWatcher.Errors:
				log.Println("error:", err)
			}
		}
	}(fsWatcher)

	if err = fsWatcher.Add("/Users/caoyuan/workstation/go-learning/practise/fsnotify-practise"); err != nil {
		log.Fatal(err)
	}

	select {}
}
