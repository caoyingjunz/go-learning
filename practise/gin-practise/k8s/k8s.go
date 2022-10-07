package k8s

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientInterface interface {
	Add(key string, obj *kubernetes.Clientset)
	Update(key string, obj *kubernetes.Clientset)
	Delete(key string)
	Get(key string) *kubernetes.Clientset
}

type ClientStore struct {
	lock sync.RWMutex

	items map[string]*kubernetes.Clientset
}

func (cs *ClientStore) Add(key string, obj *kubernetes.Clientset) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items[key] = obj
}

func (cs *ClientStore) Update(key string, obj *kubernetes.Clientset) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	cs.items[key] = obj
}

func (cs *ClientStore) Delete(key string) {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	_, exits := cs.items[key]
	if exits {
		delete(cs.items, key)
	}
}

func (cs *ClientStore) Get(key string) *kubernetes.Clientset {
	cs.lock.Lock()
	defer cs.lock.Unlock()

	item, exits := cs.items[key]
	if !exits {
		return nil
	}

	return item
}

func NewClientStore() ClientInterface {
	return &ClientStore{
		items: map[string]*kubernetes.Clientset{},
	}
}

type EngineInterface interface {
	Start(stopCh chan struct{}) error

	GetPod(ctx context.Context, key string, name string, namespace string) (*v1.Pod, error)
}

type KubeEngine struct {
	client ClientInterface

	path string
}

// Start watches for the creation and deletion of plugin sockets at the path
func (k *KubeEngine) Start(stopCh chan struct{}) error {
	fmt.Println("Plugin Watcher Start at", k.path)

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Traverse plugin dir and add filesystem watchers before starting the plugin processing goroutine.
	//if err := w.traversePluginDir(w.path); err != nil {
	//	klog.ErrorS(err, "Failed to traverse plugin socket path", "path", w.path)
	//}

	err = fsWatcher.Add(k.path)
	if err != nil {
		return err
	}

	go func(fsWatcher *fsnotify.Watcher) {
		for {
			select {
			case event := <-fsWatcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					if err = k.handleCreateEvent(event); err != nil {
						fmt.Println("create event failed", err)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if err = k.handleDeleteEvent(event); err != nil {
						fmt.Println("delete event failed", err)
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

func (k *KubeEngine) GetPod(cxt context.Context, key string, name string, namespace string) (*v1.Pod, error) {
	clientSet := k.client.Get(key)
	if clientSet == nil {
		return nil, fmt.Errorf("%s not register", key)
	}

	pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (k *KubeEngine) handleCreateEvent(event fsnotify.Event) error {
	client, err := k.newClient(event.Name)
	if err != nil {
		//打印log
		return err
	}

	fi, err := os.Stat(event.Name)
	if err != nil {
		return err
	}

	if strings.HasPrefix(fi.Name(), ".") || fi.IsDir() {
		fmt.Println("Ignoring file (starts with '.') or dir", "path", fi.Name())
		return nil
	}

	name := fi.Name()

	k.client.Add(name, client)
	fmt.Println(name, "register")
	return nil
}

func (k *KubeEngine) handleDeleteEvent(event fsnotify.Event) error {
	_, name := filepath.Split(event.Name)

	k.client.Delete(name)
	fmt.Println(name, "unregister")
	return nil
}

func (k *KubeEngine) newClient(kubeConfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func NewKubeEngine() EngineInterface {
	return &KubeEngine{
		client: NewClientStore(),
		path:   "/Users/caoyuan/workstation/go-learning/practise/gin-practise",
	}
}
