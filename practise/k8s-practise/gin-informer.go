package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// assignedPod2 selects pods that are assigned (scheduled and running).
func assignedPod2(pod *corev1.Pod) bool {
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

	sharedInformers := informers.NewSharedInformerFactory(clientset, 0)
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
				return assignedPod2(t) // 可用，可不用
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

	// 构造 pod podLister，用于 gin 的查询
	podLister := sharedInformers.Core().V1().Pods().Lister()
	// 启动 gin router
	// 仅做演示，无封装，无异常处理
	// 启动之后， curl 127.0.0.1:8088/pods 测试效果
	r := gin.Default()
	r.GET("/pods", func(c *gin.Context) {
		pod, err := podLister.Pods("default").Get("test-nginx-7b84788ff9-64fdd")
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "pong", "code": 1000, "result": pod})
	})

	_ = r.Run(":8088")
}
