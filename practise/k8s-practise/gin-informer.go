package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	sharedInformers := informers.NewSharedInformerFactory(clientSet, 0)
	// refer to https://github.com/kubernetes/kubernetes/blob/ea0764452222146c47ec826977f49d7001b0ea8c/staging/src/k8s.io/client-go/dynamic/dynamicinformer/informer_test.go#L107
	// TODO: 可以追加更多的 gvr

	gvrs := []schema.GroupVersionResource{
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "", Version: "v1", Resource: "pods"},
	}
	for _, gvr := range gvrs {
		if _, err = sharedInformers.ForResource(gvr); err != nil {
			panic(err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		pods, err := podLister.List(labels.Everything())
		if err != nil {
			panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err, "code": 400})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "pong", "code": 200, "result": pods})
	})

	_ = r.Run(":8088")
}
