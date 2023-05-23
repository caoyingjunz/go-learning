# code-generator

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
