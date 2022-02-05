# kubectl plugin 源码分析

### kubectl 概述

`kubectl` 作为 [kubernetes](https://github.com/kubernetes/kubernetes) 官方提供的命令行工具，基于 [cobra](https://github.com/spf13/cobra) 实现，用于对 `kubernetes` 集群进行管理

- 本文仅对 `kubectl plugin` 的源码进行分析，如何使用请移步 [Extend kubectl with plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- `cobra` 使用请移步 [pixiuctl](https://github.com/caoyingjunz/go-learning/tree/master/practise/cobra-practise)

### 入口函数 main
``` go
// cmd/kubectl/kubectl.go

func main() {

	command := cmd.NewDefaultKubectlCommand() // 主干函数
	...
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
```
`kubectl` 的 `main` 函数非常简洁，新建 command 之后直接执行；command 由 `NewDefaultKubectlCommand` 返回


###