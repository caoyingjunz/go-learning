package k8s

import (
	"context"
	"fmt"
	"log"
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
	err = fsWatcher.Add(k.path)
	if err != nil {
		return err
	}

	go func(fsWatcher *fsnotify.Watcher) {
		for {
			select {
			case event := <-fsWatcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					_ = k.handleCreateEvent(event)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					k.handleDeleteEvent(event)
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

	pod, err := clientSet.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (k *KubeEngine) handleCreateEvent(event fsnotify.Event) error {
	configFile := event.Name
	client, err := k.newClient(event.Name)
	if err != nil {
		//打印log
		return err
	}

	parts := strings.Split(configFile, "/")
	if len(parts) == 0 {
		return fmt.Errorf("文件路径不对")
	}

	key := parts[len(parts)-1]
	k.client.Add(key, client)

	fmt.Println(key, "register")
	return nil
}

func (k *KubeEngine) handleDeleteEvent(event fsnotify.Event) error {
	parts := strings.Split(event.Name, "/")
	if len(parts) == 0 {
		return fmt.Errorf("文件路径不对")
	}

	key := parts[len(parts)-1]
	k.client.Delete(key)

	fmt.Println(key, "unregister")
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
