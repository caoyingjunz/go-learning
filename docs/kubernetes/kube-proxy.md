# kube-proxy 源码分析

### mode: iptables

### 启动 `kube-proxy`
- 通过命令行启动 `kube-proxy`， 代码位置 `cmd/kube-proxy/proxy.go`
	```
    package main
    ...

	func main() {
		rand.Seed(time.Now().UnixNano())

		command := app.NewProxyCommand()
        ...
		if err := command.Execute(); err != nil {
			os.Exit(1)
		}
	}
	```
	- kube-proxy 使用 [cobra](https://github.com/spf13/cobra) 来新建 `ProxyCommand`, 完成配置初始化和校验，以及程序的执行

-
