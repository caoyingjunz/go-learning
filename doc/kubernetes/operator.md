# Operator

### Mac 安装 operator-sdk

```
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.4.2/operator-sdk_darwin_amd64
chmod +x operator-sdk_darwin_amd64
sudo cp operator-sdk_darwin_amd64 /usr/local/go/bin/operator-sdk

or

brew install operator-sdk
```

### 初始化项目
* [Quickstart](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/)
* [Example](http://www.dockone.io/article/8733)
* [Fulltutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

```
mkdir podset-operator
cd podset-operator
operator-sdk init --domain github.com --repo github.com/caoyingjunz/podset-operator
```

### 创建 API 和 controller

```
# Create a PodSet API with Group: cache, Version: v1beta1 and Kind: PodSet
operator-sdk create api --group cache --version v1alpha1 --kind PodSet --resource --controller
```

### 生成 CRD 然后 apply, 生成 yaml 存放在 config/crd/bases

```
make install
```

### 自定义控制器代码
TODO

### Build and push image
```
docker build -f Dockerfile . -t jacky06/podset-operator:v0.0.1
docker push jacky06/podset-operator:v0.0.1
```

### 部署 PodSet CRD

crd 来自 [podset-operator](https://github.com/caoyingjunz/podset-operator/blob/master/deploy/crds/kubez_podsets_crd.yaml)

```
kubectl apply -f deploy/crds/kubez_podsets_crd.yaml
```
