# Go Learning

## Kubernetes Plans
1. 写一个控制器，当一个应用启动的时候，保证其他的node上也有其镜像（Done）
2. Prometheus 部署(DONE) [kube-prometheus](https://github.com/coreos/kube-prometheus)
3. kubectl plugin （DONE）[plugins](https://github.com/ahmetb/kubectx)
4. 容器内部访问 k8s api （DONE）[参考](https://www.jianshu.com/p/b1a723033a3c)
5. k8s 首先看 event 和 ns (DONE)
6. k8s ingress 跨域 [ingress](https://blog.csdn.net/u012375924/article/details/94360425)
7. K8s pod 出外网 (DONE) 请求先到网关（cni），直接使用 host 的路由出外网
8. Docker 集成 ovs [ovs-docs](https://docs.openvswitch.org/en/latest/intro/install/general/#obtaining-open-vswitch-sources)
9.  学习 loadbalancer 实现 [cloud-provider-openstack](https://github.com/kubernetes/cloud-provider-openstack)
10. iptables&ipvs实现负载均衡（DONE)
11. API Watch 实现[watch](https://www.jianshu.com/p/1cb577f750f0)
12. 分析 `deployment` 创建流程（Done）
13. 分析 `pod` 创建流程（Done）
14. 分析 `scheduler` 工作流程 (Done)
15. 分析 `ListAndWatch` 原理(Done)
16. `istio` (TODO)

## Kubernetes Docs
1. [kubernetes 经典网络分析](./doc/network.md)
2. [kubectl exec 实现分析](./doc/kubernetes/kubeexec.md)
3. [kubernetes 集群快速搭建](https://github.com/caoyingjunz/kubez-ansible)
4. [自定义控制器](./doc/kubernetes/controller.md)
5. [Operator 代码分析](./doc/kubernetes/operator.md)
6. [kube-proxy代码分析](https://github.com/caoyingjunz/kubezspaces/blob/master/docs/kubernetes/kube-proxy.md)
7. [scheduler代码分析](TODO)
8. [kubelet](TODO)

## OpenStack
1. zun 代码分析 （TODO）
2. kuryr 代码分析（TODO）
3. qinling 代码分析（TODO）

## Golang
1. logrus, klog, zap 的使用 (DONE) [官方文档](https://github.com/sirupsen/logrus) | [klog官方文档](https://github.com/kubernetes/klog) | [zap官方文档](https://github.com/uber-go/zap)
2. 并发 goroutine (DONE)
3. slice 排序 (DONE)
4. channel 实现并发和锁 (DONE), 排他锁 (WIP)
5. go-restful (DONE) [官方文档](https://github.com/emicklei/go-restful)
6. cobra (DONE) [官方文档](https://github.com/spf13/cobra)
7. map 的同时读写 [官方文档](https://golang.org/pkg/sync/#Map)
8. context (DONE) [context](https://mp.weixin.qq.com/s/GpVy1eB5Cz_t-dhVC6BJNw)
9. 协程的调度机制(TODO)
10. interface (DOEN)
11. sync.Mutex (TODO)
12. gRPC (WIP)[官方文档](https://github.com/grpc/grpc-go)
13. yaml (DONE) [官方文档](https://github.com/go-yaml/yaml) | [yaml](https://www.jianshu.com/p/84499381a7da)
14. 匿名函数 (DONE)
15. go-gorm (DONE)[中文文档](http://gorm.book.jasperxu.com/) | [官方文档](https://github.com/go-gorm/gorm)
16. golang 操作 rabbitMq（WIP）
17. configor, ini 配置文件学习（WIP）[官方文档](https://github.com/jinzhu/configor)
