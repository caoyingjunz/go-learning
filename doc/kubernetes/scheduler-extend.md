# Kubernetes 调度扩展

## 概述
- 不涉及调度的基本原理
- Kubernetes 的调度代码分析可参考 [源码分析](https://github.com/caoyingjunz/go-learning/blob/master/doc/kubernetes/scheduler-start.md)

## 扩展方法
- 修改 `kube-scheduler` 的代码
- 独立调度器
- Scheduler Extender
- Scheduler Framework

### 修改 `kube-scheduler` 的代码
- 修改原生 `kube-scheduler` 的代码，编译，部署。
- 不推荐
  - 开发难度大，维护成本高。

### 独立调度器
- 自定义开发调度器，与 `kube-scheduler` 同时部署运行在 `kubernetes` 集群中，`pod` 通过 `spec.schedulerName` 指定调度程序。
- 不推荐
  - 开发难度大，维护成本高，多个调度器之间存在资源冲突问题。

### Scheduler Extender
通过 webhook 实现 filter 和 score
TODO

### Scheduler Framework
TODO