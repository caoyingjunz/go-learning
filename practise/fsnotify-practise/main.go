package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// more info https://github.com/fsnotify/fsnotify
// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/pluginmanager/pluginwatcher/plugin_watcher.go#L60

func handleCreateEvent(event fsnotify.Event) error {
	log.Println("Add file", event)
	fi, err := os.Stat(event.Name)
	if err != nil {
		return err
	}

	fmt.Println(fi.Name())
	return nil
}

func handleDeleteEvent(event fsnotify.Event) error {
	log.Println("Delete file", event)

	_, name := filepath.Split(event.Name)
	fmt.Println(name)
	return nil
}

func Start(stopCh <-chan struct{}) error {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// Traverse plugin dir and add filesystem watchers before starting the plugin processing goroutine.
	//if err := w.traversePluginDir(w.path); err != nil {
	//	klog.ErrorS(err, "Failed to traverse plugin socket path", "path", w.path)
	//}

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
					err = handleCreateEvent(event)
					if err != nil {
						fmt.Println("create event failed", err)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					err = handleDeleteEvent(event)
					if err != nil {
						fmt.Println("remove event failed", err)
					}
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
