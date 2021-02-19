# Informer 原理分析

#### 开局一张图
![informer](../images/informer.png)

#### client-go components
* Reflector:A reflector, which is defined in type Reflector inside package cache, watches the Kubernetes API for the specified resource type (kind). The function in which this is done is ListAndWatch. The watch could be for an in-built resource or it could be for a custom resource. When the reflector receives notification about existence of new resource instance through the watch API, it gets the newly created object using the corresponding listing API and puts it in the Delta Fifo queue inside the watchHandler function.

* Informer: An informer defined in the base controller inside package cache pops objects from the Delta Fifo queue. The function in which this is done is processLoop. The job of this base controller is to save the object for later retrieval, and to invoke our controller passing it the object.

* Indexer: An indexer provides indexing functionality over objects. It is defined in type Indexer inside package cache. A typical indexing use-case is to create an index based on object labels. Indexer can maintain indexes based on several indexing functions. Indexer uses a thread-safe data store to store objects and their keys. There is a default function named MetaNamespaceKeyFunc defined in type Store inside package cache that generates an object’s key as <namespace>/<name> combination for that object.

#### 姗姗来迟的 Demo
    package main

    import (
        "time"

        v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/informers"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/tools/cache"
        "k8s.io/client-go/tools/clientcmd"
    )

    func main() {
        // build kubeconfig
        config, err := clientcmd.BuildConfigFromFlags("", "kubeconfigPath")
        if err != nil {
            panic(err)
        }

        // New clientset by kubeconfig
        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
            panic(err)
        }

        stopCh := make(chan struct{})
        defer close(stopCh)

        // 构造 sharedInformers，每分钟同步一次
        sharedInformers := informers.NewSharedInformerFactory(clientset, time.Minute)

        informer := sharedInformers.Core().V1().Services().Informer()

        // 新建 EventHandler，需要实现 3 种回调方法: AddFunc, UpdateFunc, DeleteFunc
        informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                ...
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                ...
            },
            DeleteFunc: func(obj interface{}) {
                ...
            },
        })

        // Run informer
        informer.Run(stopCh)
    }
- 完整 `demo` 请参考 [informer-demo](.../../../../examples/informer-demo.go)

#### Informer 代码分析
- TODO
- TODO
