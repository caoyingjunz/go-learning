# CSI Plugin 注册机制源码分析

### node-driver-registrar 流程解析：

1. rpc 调用 nfs-csi-plugin 中的 identity server 中的 GetPluginInfo 方法，返回我们自研plugin相关的基本信息， 获取 csi driver name

2. 启动 grpc server，并监听在宿主机的 /var/lib/kubelet/plugins_registry/${csiDriverName}-reg.sock （对应容器的/registration/${driver-name}-reg.sock）

3. 该 grpc server 提供 GetInfo 和 NotifyRegistrationStatus 方法供 kubelet plugin manager 调用，


### kubelet plugin manager 流程解析

1. 当 ds 部署的 node-driver-registrar sidecar 启动时，/var/lib/kubelet/plugins_registry 会新增一个 socket file nfs.csi.k8s.io-reg.sock

2. 这个 nfs.csi.k8s.io-reg.sock 会被 plugin watcher 监听并写入desiredStateOfWorld缓存中（通过github.com/fsnotify/fsnotify实现）

3. reconciler 会对比缓存，找出需要新增或者删除的 csi plugin。当新增时，reconciler 会通过 grpc client（/var/lib/kubelet/plugins_registry）调用 node-driver-registrar sidecar container中rpc server提供的的GetInfo，然后根据返回字段的type，找到对应的 plugin hander，

```go
// socketPath 为 /registration/${driver-name}-reg.sock
client, conn, err := dial(socketPath, dialTimeoutDuration)
        // 调用node-driver-registrar sidecar container中rpc server提供的的GetInfo
        infoResp, err := client.GetInfo(ctx, &registerapi.InfoRequest{})
        // 这里handler就是上文说的CSIPlugin type的csi.RegistrationHandler{}对象
        handler, ok := pluginHandlers[infoResp.Type]
        // 调用handler.ValidatePlugin
        if err := handler.ValidatePlugin(infoResp.Name, infoResp.Endpoint, infoResp.SupportedVersions); err != nil {
        }
        ...
        // 加入actualStateOfWorldUpdater缓存
        err = actualStateOfWorldUpdater.AddPlugin(cache.PluginInfo{
            SocketPath: socketPath,
            Timestamp:  timestamp,
            Handler:    handler,
            Name:       infoResp.Name,
        })
// infoResp.Endpoint 就是 nfs-csi 的 socket path
if err := handler.RegisterPlugin(infoResp.Name, infoResp.Endpoint, infoResp.SupportedVersions); err != nil {
            return og.notifyPlugin(client, false, fmt.Sprintf("RegisterPlugin error -- plugin registration failed with err: %v", err))
        }
    ...
    return registerPluginFunc
```

4. plugin handler 会根据传入的 csi-plugin 监听的 socket path，直接和我们 nfs-csi-plugin 通信，并调用该对象的 ValidatePlugin 和 RegisterPlugin 来注册插件，这里的注册插件其实就是设置 node annotation 和 创建/更新CSINode对象。

```bash
metadata:
 annotations:
   csi.volume.kubernetes.io/nodeid: '{"nfs.csi.k8s.io":"kube-master"}'

# kubectl get csinode
NAME           DRIVERS   AGE
kubez-master   1         101d
kubez-node1    1         99d
kubez-node2    1         101d
```

5. kubelet plugin-magner 通过 rpc 调用 NotifyRegistrationStatus 告知 node-driver-registrar 注册结果。
