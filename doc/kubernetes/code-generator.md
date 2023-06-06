# code-generator

- deepcopy-gen: 生成深度拷贝方法，为每个 T 类型生成 func (t* T) DeepCopy() *T 方法，API 类型都需要实现深拷贝
- client-gen: 为资源生成标准的 clientset
- informer-gen: 生成 informer，提供事件机制来响应资源的事件
- lister-gen: 生成 Lister，为 get 和 list 请求提供只读缓存层（通过 indexer 获取）

```shell
git clone https://github.com/kubernetes/code-generator.git
git checkout v0.24.9
```

```shell
# 进行安装
go install ./cmd/{client-gen,deepcopy-gen,informer-gen,lister-gen}
```

```shell
# 验证安装
client-gen -h
```

```shell
代码结构可以参考 
https://github.com/kubernetes/sample-controller
```

```shell
$ mkdir test && cd test
$ go mod init test
$ mkdir -p pkg/apis/example.com/v1
➜ test tree
├── go.mod
├── go.sum
└── pkg
    └── apis
        └── example.com
            └── v1
                ├── doc.go
                ├── register.go
                └── types.go
```
