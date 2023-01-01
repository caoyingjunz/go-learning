package main

import (
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/controller-manager/pkg/clientbuilder"
)

const ResourceResyncTime time.Duration = 0

// monitor runs a Controller with a local stop channel.
type monitor struct {
	controller cache.Controller
	store      cache.Store

	stopCh chan struct{}
}

func (m *monitor) Run() {
	// 启动监控
	m.controller.Run(m.stopCh)
}

var (
	sharedInformers informers.SharedInformerFactory
	restMapper      *restmapper.DeferredDiscoveryRESTMapper
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	rootClientBuilder := clientbuilder.SimpleControllerClientBuilder{ClientConfig: config}

	versionedClient := rootClientBuilder.ClientOrDie("shared-informers")
	sharedInformers = informers.NewSharedInformerFactory(versionedClient, ResourceResyncTime)

	// Use a discovery client capable of being refreshed.
	discoveryClient := rootClientBuilder.DiscoveryClientOrDie("controller-discovery")
	cachedClient := cacheddiscovery.NewMemCacheClient(discoveryClient)
	restMapper = restmapper.NewDeferredDiscoveryRESTMapper(cachedClient)

	gvrs := []schema.GroupVersionResource{
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "", Version: "v1", Resource: "pods"},
	}

	monitors := make(map[schema.GroupVersionResource]*monitor)
	for _, gvr := range gvrs {
		kind, err := restMapper.KindFor(gvr)
		if err != nil {
			panic(err)
		}

		c, s, err := controllerFor(gvr, kind)
		if err != nil {
			panic(err)
		}
		monitors[gvr] = &monitor{store: s, controller: c}
	}

	stopCh := make(chan struct{})
	for kind, monitor := range monitors {
		if monitor.stopCh == nil {
			monitor.stopCh = make(chan struct{})
			sharedInformers.Start(stopCh)
			fmt.Println("start monitor for ", kind)
			go monitor.Run()
		}
	}

	select {}
}

func controllerFor(resource schema.GroupVersionResource, kind schema.GroupVersionKind) (cache.Controller, cache.Store, error) {
	handlers := cache.ResourceEventHandlerFuncs{
		// add the event to the dependencyGraphBuilder's graphChanges.
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			fmt.Println("add", kind.String(), mObj.GetName(), mObj.GetFinalizers())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			mObj := newObj.(v1.Object)
			fmt.Println("update", kind.String(), mObj.GetName(), mObj.GetFinalizers())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			fmt.Println("delete", kind.String(), mObj.GetName(), mObj.GetFinalizers())
		},
	}
	shared, err := sharedInformers.ForResource(resource)
	if err != nil {
		panic(err)
	}
	shared.Informer().AddEventHandlerWithResyncPeriod(handlers, ResourceResyncTime)

	return shared.Informer().GetController(), shared.Informer().GetStore(), nil
}
