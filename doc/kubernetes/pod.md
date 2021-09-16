### Pod 代码流程分析

1.  `kube-controller-manager` 在启动的时候会初始化全部 `kubernetes` 控制器， 其中包含管理 `pod` 的 `podgc` 控制器

```
// NewControllerManagerCommand 创建一个cobra的命令对象
func NewControllerManagerCommand() *cobra.Command {
    ...
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
            ...
			c, err := s.Config(KnownControllers(), ControllersDisabledByDefault.List())
			...
		},
	}

# 初始化 startPodGCController
func NewControllerInitializers(loopMode ControllerLoopMode) map[string]InitFunc {
	...
	controllers["podgc"] = startPodGCController
    ...

# 启动 pod 控制器
func startPodGCController(ctx ControllerContext) (http.Handler, bool, error) {
	go podgc.NewPodGC(
		ctx.ClientBuilder.ClientOrDie("pod-garbage-collector"),
		ctx.InformerFactory.Core().V1().Pods(),
		int(ctx.ComponentConfig.PodGCController.TerminatedPodGCThreshold),
	).Run(ctx.Stop)
	return nil, true, nil
}
```

2.
