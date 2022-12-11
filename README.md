# Go-Learning Overview

Go-Learning's mission statement is:

    To provide a learning and practise map for OpenStack, kubernetes, golang and the others.

go-learning 适用于有一定 `kubernetes` 经验，且想更进一步的同学。

- kubernetes 功能增强 [Pixiu(貔貅)](https://github.com/caoyingjunz/pixiu)
- 快速部署 [kubez-ansible](https://github.com/caoyingjunz/kubez-ansible)
- workload 自动扩缩容 [piuxiu-autoscaler](https://github.com/caoyingjunz/pixiu-autoscaler)

## Kubernetes
- [kubectl plugin 源码分析](./doc/kubernetes/kubectl-plugin.md)
- [scheduler 源码分析一 - 启动](./doc/kubernetes/scheduler-start.md)
- [scheduler 源码分析二 - 调度](./doc/kubernetes/scheduler-schedule.md)
- [kubernetes 网络分析](./doc/kubernetes/network.md)
- [kube-proxy 源码分析](./doc/kubernetes/kube-proxy.md)
- [operator 用法详解](./doc/kubernetes/operator.md)
- [CSI 注册机制源码分析](./doc/kubernetes/csi.md)
- [cloud-provider-openstack](https://github.com/kubernetes/cloud-provider-openstack)

## Examples
- [Examples](./examples/README.md) 提供丰富的 `kubernetes` 用法举例.
- [pixiuctl](https://github.com/caoyingjunz/go-learning/tree/master/practise/cobra-practise/pixiuctl) 基于 [cobra](https://github.com/spf13/cobra) 实现命令行
  - subcommand
  - plugin
- [gRPC Usage](./practise/grpc-practise/README.md)
- [gin&informer](./practise/k8s-practise/gin-informer.go) 提供 `gin` 调用 `informer` 的用法
- [webShell](https://github.com/caoyingjunz/kube-webshell) 提供 `webshell` 调用 `kubernetes` `pod` 的用法展示

## TODO
- scheduler 代码分析(WIP)
- kubelet 代码分析
- 微服务学习（istio）
- gc 机制分析
- pod 驱逐代码分析

Copyright 2019 caoyingjun (cao.yingjunz@gmail.com) Apache License 2.0
