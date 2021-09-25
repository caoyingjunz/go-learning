package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// more info https://github.com/fsnotify/fsnotify
// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/pluginmanager/pluginwatcher/plugin_watcher.go#L60

func handleCreateEvent(event fsnotify.Event) error {
	log.Println("Add file", event)
	return nil
}

func handleDeleteEvent(event fsnotify.Event) error {
	log.Println("Delete file", event)
	return nil
}

func Start(stopCh <-chan struct{}) error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// 监听两个文件夹
	if err = fsWatcher.Add("/Users/caoyuan/workstation/go-learning/practise/fsnotify-practise"); err != nil {
		log.Fatal(err)
	}
	if err = fsWatcher.Add("/Users/caoyuan/workstation/go-learning"); err != nil {
		log.Fatal(err)
	}

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
			case <-stopCh:
				fsWatcher.Close()
			}
		}
	}(fsWatcher)

	return nil
}

func main() {
	stopCh := make(chan struct{})

	if err := Start(stopCh); err != nil {
		log.Fatal(err)
	}

	select {}
}
