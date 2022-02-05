# kube-proxy 源码分析

### mode
- iptables

### 启动 `kube-proxy`
- 通过命令行启动 `kube-proxy`， 代码位置 `cmd/kube-proxy/proxy.go`
	``` go
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

- `cobra` 在调用 `command.Execute` 的时候会运行一个指定的 `ProxyServer`，并运行 `runLoop`.
	``` go
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
	``` go
	func newProxyServer(
		config *proxyconfigapi.KubeProxyConfiguration,
		cleanupAndExit bool,
		master string) (*ProxyServer, error) {
        ...
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
        - `mode`:
			- `iptables` 或者 `""(不填)`: `iptables Proxier`, 调用 `iptables.NewProxier`
			- `ipvs`: `ipvs Proxier`, 调用 `ipvs.NewProxier`

- 本文仅分析 `mode` 为 `iptables` 场景;  `NewProxyServer` 会调用 `iptables.NewProxier` 方法来初始化一个 `proxier`
	``` go
	func NewProxier(...) (*Proxier, error) {
		// 设置 route_localnet = 1
		if val, _ := sysctl.GetSysctl(sysctlRouteLocalnet); val != 1 {
			if err := sysctl.SetSysctl(sysctlRouteLocalnet, 1); err != nil {
				return nil, fmt.Errorf("can't set sysctl %s: %v", sysctlRouteLocalnet, err)
			}
		}

		// 确保 br_netfilter 和 bridge-nf-call-iptables 被开启, container 连接到linuxbridge的时，需要两者被开启.
		if val, err := sysctl.GetSysctl(sysctlBridgeCallIPTables); err == nil && val != 1 {
			klog.Warning("missing br-netfilter module or unset sysctl br-nf-call-iptables; proxy may not work as intended")
		}

		// 为 SNAT iptables 规则生成 masquerade 标记
		masqueradeValue := 1 << uint(masqueradeBit)
		masqueradeMark := fmt.Sprintf("%#08x", masqueradeValue)
		...

		// 初始化 proxier
		proxier := &Proxier{
			portsMap:                 make(map[utilproxy.LocalPort]utilproxy.Closeable),
			serviceMap:               make(proxy.ServiceMap),
			serviceChanges:           proxy.NewServiceChangeTracker(newServiceInfo, &isIPv6, recorder),
			endpointsMap:             make(proxy.EndpointsMap),
			endpointsChanges:         proxy.NewEndpointChangeTracker(hostname, newEndpointInfo, &isIPv6, recorder, endpointSlicesEnabled),
			syncPeriod:               syncPeriod,
			iptables:                 ipt,
			masqueradeAll:            masqueradeAll,
			masqueradeMark:           masqueradeMark,
			exec:                     exec,
			localDetector:            localDetector,
			hostname:                 hostname,
			nodeIP:                   nodeIP,
			portMapper:               &listenPortOpener{},
			recorder:                 recorder,
			serviceHealthServer:      serviceHealthServer,
			healthzServer:            healthzServer,
			precomputedProbabilities: make([]string, 0, 1001),
			iptablesData:             bytes.NewBuffer(nil),
			existingFilterChainsData: bytes.NewBuffer(nil),
			filterChains:             bytes.NewBuffer(nil),
			filterRules:              bytes.NewBuffer(nil),
			natChains:                bytes.NewBuffer(nil),
			natRules:                 bytes.NewBuffer(nil),
			nodePortAddresses:        nodePortAddresses,
			networkInterfacer:        utilproxy.RealNetwork{},
		}

		// 初始化 syncRunner, 设置 proxier.syncProxyRules 方法作为一个参数构造 proxier.syncRunner
		proxier.syncRunner = async.NewBoundedFrequencyRunner("sync-runner", proxier.syncProxyRules, minSyncPeriod, time.Hour, burstSyncs)

		// 启动 ipt.Monitor
		go ipt.Monitor(utiliptables.Chain("KUBE-PROXY-CANARY"),
			[]utiliptables.Table{utiliptables.TableMangle, utiliptables.TableNAT, utiliptables.TableFilter},
			proxier.syncProxyRules, syncPeriod, wait.NeverStop)
		return proxier, nil
	}
	```
	- `NewProxier` 方法主要完成如几件事:
        - 设置 `route_localnet` = 1
        - 检查, 确保 `br_netfilter` 和 `bridge-nf-call-iptables` = 1
        - 为 `SNAT` `iptables` 规则生成 `masquerade` 标记
        - 初始化 `proxier`
        - 初始化 `syncRunner`, 设置 `proxier.syncProxyRules` 方法作为参数构造 `syncRunner`
        - 启动一个 `goroutine`，用于启动 `ipt.Monitor`

- 完成 `Proxier` 创建之后, `Run` 方法会调用 `o.runLoop`，通过 goroutine 启动 `o.proxyServer.Run`, 代码位置 `cmd/kube-proxy/app/server.go`
	``` go
	func (o *Options) runLoop() error {
		...
		// 通过 goroutine 启动 proxy
		go func() {
			err := o.proxyServer.Run()
			o.errCh <- err
		}()
		...
	}
	```

- `proxyServer.Run`
	``` go
	// This should never exit (unless CleanupAndExit is set).
	func (s *ProxyServer) Run() error {
		...

		// Start up a metrics server if requested
		if len(s.MetricsBindAddress) > 0 {
			...
		}

		// Tune conntrack, if requested. Conntracker is always nil for windows
		if s.Conntracker != nil {
			...
		}
		...

		// Make informers that filter out objects that want a non-default service proxy.
		informerFactory := informers.NewSharedInformerFactoryWithOptions(s.Client, s.ConfigSyncPeriod,
			informers.WithTweakListOptions(func(options *metav1.ListOptions) {
				options.LabelSelector = labelSelector.String()
			}))

		// Create configs (i.e. Watches for Services and Endpoints or EndpointSlices)
		// Note: RegisterHandler() calls need to happen before creation of Sources because sources
		// only notify on changes, and the initial update (on process start) may be lost if no handlers
		// are registered yet.
		serviceConfig := config.NewServiceConfig(informerFactory.Core().V1().Services(), s.ConfigSyncPeriod)
		serviceConfig.RegisterEventHandler(s.Proxier)
		go serviceConfig.Run(wait.NeverStop)

		if s.UseEndpointSlices {
			endpointSliceConfig := config.NewEndpointSliceConfig(informerFactory.Discovery().V1beta1().EndpointSlices(), s.ConfigSyncPeriod)
			endpointSliceConfig.RegisterEventHandler(s.Proxier)
			go endpointSliceConfig.Run(wait.NeverStop)
		} else {
			endpointsConfig := config.NewEndpointsConfig(informerFactory.Core().V1().Endpoints(), s.ConfigSyncPeriod)
			endpointsConfig.RegisterEventHandler(s.Proxier)
			go endpointsConfig.Run(wait.NeverStop)
		}

		// This has to start after the calls to NewServiceConfig and NewEndpointsConfig because those
		// functions must configure their shared informer event handlers first.
		informerFactory.Start(wait.NeverStop)

		...

		// Just loop forever for now...
		s.Proxier.SyncLoop()
		return nil
	}
	```
	- `s.Run` 方法主要完成:
        - Start up a metrics server if requested
        - Tune conntrack, if requested
        - informerFactory.Start, more info: [informer](./informer.md)
        - 启动 s.Proxier.SyncLoop 方法

- 在 `s.Run` 中调用 `SyncLoop` 方法 进行 `Loop`
	``` go
	func (proxier *Proxier) SyncLoop() {
		...
		proxier.syncRunner.Loop(wait.NeverStop)
	}

	func (bfr *BoundedFrequencyRunner) Loop(stop <-chan struct{}) {
		klog.V(3).Infof("%s Loop running", bfr.name)
		bfr.timer.Reset(bfr.maxInterval)
		for {
			select {
			...
			case <-bfr.timer.C():
				bfr.tryRun()
			case <-bfr.run:
				bfr.tryRun()
			case <-bfr.retry:
				bfr.doRetry()
			}
		}
	}

	func (bfr *BoundedFrequencyRunner) tryRun() {
		bfr.mu.Lock()
		defer bfr.mu.Unlock()

		if bfr.limiter.TryAccept() {
			...
			bfr.fn()  // 真正执行的方法，调用 fn(), 根据初始化传入的方法参数，fn = syncProxyRules
			...
			return
		}

		...
	}
	```
	- `SyncLoop` 循环的调用 `syncProxyRules`，实现对 `service` 和 `ingress` 的 `iptables` 规则下发

- 循环运行核心方法 `syncProxyRules`, 完成 `kube-proxy` 职能. 代码位置 `pkg/proxy/iptables/proxier.go`.

	``` go
	func (proxier *Proxier) syncProxyRules() {
		proxier.mu.Lock()
		defer proxier.mu.Unlock()
		...

		TODO: 方法特别长（900+行), 后续补充
	}
	```
    - `syncProxyRules` 方法主要完成:
        - 感知到 services 和 endpoints 的 changed, 然后完成 iptables 规则的下发
