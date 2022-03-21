package main

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// assignedPod selects pods that are assigned (scheduled and running).
func assignedPod(pod *corev1.Pod) bool {
	return pod.Spec.NodeName == "test"
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/caoyuan/.kube/config")
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stopCh := server.SetupSignalHandler()
		<-stopCh
		cancel()
	}()

	sharedInformers := informers.NewSharedInformerFactory(clientset, 0)

	// 用法 1
	//informer := sharedInformers.Core().V1().Services().Informer()
	//informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	//	AddFunc: func(obj interface{}) {
	//		mObj := obj.(v1.Object)
	//		log.Printf("New Services Added to Store: %s", mObj.GetName())
	//	},
	//	UpdateFunc: func(oldObj, newObj interface{}) {
	//		oObj := oldObj.(v1.Object)
	//		nObj := newObj.(v1.Object)
	//		log.Printf("%s Services Updated to %s", oObj.GetName(), nObj.GetName())
	//	},
	//	DeleteFunc: func(obj interface{}) {
	//		mObj := obj.(v1.Object)
	//		log.Printf("Services Deleted from Store: %s", mObj.GetName())
	//	},
	//})
	//informer.Run(ctx.Done())

	sharedInformers.Core().V1().Services().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("New Service Added to Store: %s", mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nObj := newObj.(v1.Object)
			log.Printf("%s Service Updated to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("Service Deleted from Store: %s", mObj.GetName())
		},
	})

	sharedInformers.Core().V1().Pods().Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				mObj := obj.(v1.Object)
				log.Printf("New Pod Added to Store: %s", mObj.GetName())
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oObj := oldObj.(v1.Object)
				nObj := newObj.(v1.Object)
				log.Printf("%s Pod Updated to %s", oObj.GetName(), nObj.GetName())
			},
			DeleteFunc: func(obj interface{}) {
				mObj := obj.(v1.Object)
				log.Printf("Pod Deleted from Store: %s", mObj.GetName())
			},
		},
		FilterFunc: func(obj interface{}) bool {
			switch t := obj.(type) {
			case *corev1.Pod:
				return assignedPod(t)
			case cache.DeletedFinalStateUnknown:
				if _, ok := t.Obj.(*corev1.Pod); ok {
					// The carried object may be stale, so we don't use it to check if
					// it's assigned or not. Attempting to cleanup anyways.
					return true
				}
				log.Printf("handle DeletedFinalStateUnknown error")
				return false
			default:
				log.Printf("handle object error")
				return false
			}
		},
	})

	// Start all informers.
	sharedInformers.Start(ctx.Done())
	// Wait for all caches to sync.
	sharedInformers.WaitForCacheSync(ctx.Done())

	log.Printf("informers has been started")
	select {}
}
