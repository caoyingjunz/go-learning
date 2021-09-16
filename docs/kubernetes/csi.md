# CSI Plugin 注册机制源码分析

### node-driver-registrar 流程解析

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

### external-provisioner 流程解析

1. external-provisioner 是运行在 k8s 集群内部的 controller，它作为 sidecar container 和 csi-nfs 运行在同一个 pod 中.

```bash
containers:
- args:
 - -v=2
 - --csi-address=$(ADDRESS)
 - --leader-election
 env:
 - name: ADDRESS
   value: /csi/csi.sock
 image: jacky06/csi-provisioner:v2.1.0
 name: csi-provisioner
 volumeMounts:
 - mountPath: /csi
   name: socket-dir
- args:
 - -v=5
 - --nodeid=$(NODE_ID)
 - --endpoint=$(CSI_ENDPOINT)
 env:
 - name: NODE_ID
   valueFrom:
     fieldRef:
       apiVersion: v1
       fieldPath: spec.nodeName
 - name: CSI_ENDPOINT
   value: unix:///csi/csi.sock
 image: mcr.microsoft.com/k8s/csi/nfs-csi:latest
```

2. external-provisioner 通过 informer 机制去监听 pvc/pv，当新建 pv 或者 pvc 的时候，external-provisioner 会通过 grpc (/csi/csi.sock) 调用 nfs 的 CreateVolume(DeleteVolume)方法来实际创建一个外部存储 volume

```go
// NewCSIProvisioner creates new CSI provisioner.
//
// vaLister is optional and only needed when VolumeAttachments are
// meant to be checked before deleting a volume.
func NewCSIProvisioner(client kubernetes.Interface,
	connectionTimeout time.Duration,
	...
) controller.Provisioner {
    ...
	csiClient := csi.NewControllerClient(grpcClient)

	provisioner := &csiProvisioner{
		client:                                client,
		grpcClient:                            grpcClient,
		csiClient:                             csiClient,
		snapshotClient:                        snapshotClient,
		...
		driverName:                            driverName,
		pluginCapabilities:                    pluginCapabilities,
		controllerCapabilities:                controllerCapabilities,
        ...
		scLister:                              scLister,
		csiNodeLister:                         csiNodeLister,
		nodeLister:                            nodeLister,
		claimLister:                           claimLister,
		vaLister:                              vaLister,
		eventRecorder:                         eventRecorder,
	}
    ...
	return provisioner
}

// 创建
func (p *csiProvisioner) Provision(ctx context.Context, options controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
    ...
	createCtx := markAsMigrated(ctx, result.migratedVolume)
	createCtx, cancel := context.WithTimeout(createCtx, p.timeout)
	defer cancel()
	rep, err := p.csiClient.CreateVolume(createCtx, req)
    ...
}

// 删除
func (p *csiProvisioner) Delete(ctx context.Context, volume *v1.PersistentVolume) error {
    ...
	_, err = p.csiClient.DeleteVolume(deleteCtx, &req)
    ...
```

### kubelet 通过调用 csiClient 去完成 pvc 的挂载

```go
type csiClient interface {
        NodeGetInfo(ctx context.Context) (
                nodeID string,
                maxVolumePerNode int64,
                accessibleTopology map[string]string,
                err error)
        NodePublishVolume(
                ctx context.Context,
                volumeid string,
                readOnly bool,
                stagingTargetPath string,
                targetPath string,
                accessMode api.PersistentVolumeAccessMode,
                publishContext map[string]string,
                volumeContext map[string]string,
                secrets map[string]string,
                fsType string,
                mountOptions []string,
        ) error
        NodeExpandVolume(ctx context.Context, volumeid, volumePath string, newSize resource.Quantity) (resource.Quantity, error)
        NodeUnpublishVolume(
                ctx context.Context,
                volID string,
                targetPath string,
        ) error
        NodeStageVolume(ctx context.Context,
                volID string,
                publishVolumeInfo map[string]string,
                stagingTargetPath string,
                fsType string,
                accessMode api.PersistentVolumeAccessMode,
                secrets map[string]string,
                volumeContext map[string]string,
                mountOptions []string,
        ) error
...
```
