# Informer 原理分析

#### 开局一张图
![informer](../images/informer.png)

#### client-go components
* Reflector: 用来直接和 kuber api server 通信，内部实现了 listwatch 机制，listwatch 就是用来监听资源变化的，一个listwatch 只对应一个资源，这个资源可以是 k8s 中内部的资源也可以是自定义的资源，当收到资源变化时(创建、删除、修改)时会将资源放到 Delta Fifo 队列中.

* Informer: 监听的资源的一个代码抽象，在 controller 的驱动下运行，能够将 delta filo 队列中的数据弹出处理.

* Indexer: 构建安全的本地储存，在自定义 controller 中处理对象时就是基于对象的索引在本地缓存将对象查询出来进行处理.

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
