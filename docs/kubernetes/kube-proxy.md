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

- 调用 `NewProxyServer` 新建一个 `ProxyServer`
	```
	func newProxyServer(
		config *proxyconfigapi.KubeProxyConfiguration,
		cleanupAndExit bool,
		master string) (*ProxyServer, error) {
		'''

		// 创建一个 iptables 的 utils
		execer := exec.New()
		...

		// 创建 k8s clientSet 和 eventClient
		client, eventClient, err := createClients(configClientConnection, master)
		if err != nil {
			return nil, err
		}
		...

		// 从配置文件获取代理模式：userspace，iptables，ipvs，默认为 iptables
		proxyMode := getProxyMode(string(config.Mode), kernelHandler, ipsetInterface, iptables.LinuxKernelCompatTester{})
		...

		// proxy mode 为 iptables
		if proxyMode == proxyModeIPTables {
			klog.V(0).Info("Using iptables Proxier.")
			if config.IPTables.MasqueradeBit == nil {
				// MasqueradeBit must be specified or defaulted.
				return nil, fmt.Errorf("unable to read IPTables MasqueradeBit from config")
			}

			// 判断是否开启 ipv6 双栈
			if utilfeature.DefaultFeatureGate.Enabled(features.IPv6DualStack) {
				...
			} else {
				var localDetector proxyutiliptables.LocalTrafficDetector
				localDetector, err = getLocalDetector(detectLocalMode, config, iptInterface, nodeInfo)
				if err != nil {
					return nil, fmt.Errorf("unable to create proxier: %v", err)
				}

				// TODO this has side effects that should only happen when Run() is invoked.
				proxier, err = iptables.NewProxier(
					iptInterface,
					utilsysctl.New(),
					execer,
					config.IPTables.SyncPeriod.Duration,
					config.IPTables.MinSyncPeriod.Duration,
					config.IPTables.MasqueradeAll,
					int(*config.IPTables.MasqueradeBit),
					localDetector,
					hostname,
					nodeIP,
					recorder,
					healthzServer,
					config.NodePortAddresses,
				)
			}
			...

			// 返回 ProxyServer 的实例
			return &ProxyServer{
			Client:                 client,
			EventClient:            eventClient,
			IptInterface:           iptInterface,
			IpvsInterface:          ipvsInterface,
			IpsetInterface:         ipsetInterface,
			execer:                 execer,
			Proxier:                proxier,
			Broadcaster:            eventBroadcaster,
			Recorder:               recorder,
			...
		}, nil
	}
	```
	- `NewProxyServer` 方法会根据 `mode` 来判断所使用的 `Proxier`; 默认为 `iptables`.


