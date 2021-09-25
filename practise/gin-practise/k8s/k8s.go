package k8s

import (
	"sync"

	"k8s.io/client-go/kubernetes"
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
}

type KubeEngine struct {
	client ClientInterface
}

func (k *KubeEngine) Start(stopCh chan struct{}) error {

	return nil
}

func NewKubeEngine() EngineInterface {
	return &KubeEngine{
		client: NewClientStore(),
	}
}
