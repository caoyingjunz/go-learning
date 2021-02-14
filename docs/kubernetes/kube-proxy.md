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
	- kube-proxy 使用 [cobra](https://github.com/spf13/cobra) 来新建 `NewProxyCommand`, 完成配置的初始化和校验，以及程序的执行， cobra 的用法因为篇幅有限，需自行学习.

- cobra 在调用 `command.Execute` 的时候会运行一个指定的 `ProxyServer`，并运行 `runLoop`.
	```
	func (o *Options) Run() error {
		defer close(o.errCh)
        ...
		proxyServer, err := NewProxyServer(o)
		if err != nil {
			return err
		}
		...
		o.proxyServer = proxyServer
		return o.runLoop()
	}
	```

- 调用 `NewProxyServer` 初始化 `iptables` 代理
