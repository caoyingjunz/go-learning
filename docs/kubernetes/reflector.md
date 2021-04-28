# Reflector

![informer](../images/informer.png)

``` go
// 代码来源于 kubelet.go +523
serviceIndexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
if kubeDeps.KubeClient != nil {
    serviceLW := cache.NewListWatchFromClient(kubeDeps.KubeClient.CoreV1().RESTClient(), "services", metav1.NamespaceAll, fields.Everything())
    r := cache.NewReflector(serviceLW, &v1.Service{}, serviceIndexer, 0)
    go r.Run(wait.NeverStop)
}
serviceLister := corelisters.NewServiceLister(serviceIndexer)
````

- 新建 `Indexer`: 安全的本地存储，用于存储 `Reflector` 通过 `ListWatch` 获取到的对象, 并提供获取对象的索引 <namespace>/<name>
- 构建 `Reflector`: 封装 `ListWatch` 接口和 `Indexer`, `ListWatch` 从 `api` 中获取对象，在保存到 `Indexer`
- `serviceLister`: 提供索引获取 `kubernetes` 对象的能力
