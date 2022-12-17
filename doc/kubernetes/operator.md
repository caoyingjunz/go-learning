# Operator 用法详解

### Mac 安装 operator-sdk
``` bash
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.21.0/operator-sdk_darwin_amd64
chmod +x operator-sdk_darwin_amd64
sudo cp operator-sdk_darwin_amd64 /usr/local/go/bin/operator-sdk

or

brew install operator-sdk
```

### 文档
* [官方文档](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/)
* [Example](http://www.dockone.io/article/8733)

### 版本要求
- golang 1.17
- operator-sdk v1.21.0

### 初始化 operator 项目
``` bash
mkdir podset-operator
cd podset-operator
operator-sdk init --domain pixiu.io --repo github.com/caoyingjunz/podset-operator
```

### 创建 API 和 controller
``` bash
# Create a PodSet API with Group: pixiu, Version: v1beta1 and Kind: PodSet
operator-sdk create api --group pixiu --version v1alpha1 --kind PodSet --resource --controller
```

### 创建 Webhook
``` bash
operator-sdk create webhook --group pixiu --version v1alpha1 --kind PodSet --defaulting --programmatic-validation
```

### 生成 CRD, 生成 yaml 存放在 config/crd/bases
``` bash
make manifests
```

### 自定义控制器代码
```go
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PodSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("request", req)
	log.Info("reconciling operator")

	podSet := &pixiuv1alpha1.PodSet{}
	if err := r.Get(ctx, req.NamespacedName, podSet); err != nil {
		if apierrors.IsNotFound(err) {
			// Req object not found, Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		} else {
			log.Error(err, "error requesting pod set operator")
			// Error reading the object - requeue the request.
			return reconcile.Result{Requeue: true}, nil
		}
	}
    ...
```

### Build and push image
``` bash
docker build -f Dockerfile . -t jacky06/podset-operator:v0.0.1
docker push jacky06/podset-operator:v0.0.1
```

### 部署 PodSet CRD
CRD 来自 [podset-operator](https://github.com/caoyingjunz/podset-operator/tree/master/config/crd)
