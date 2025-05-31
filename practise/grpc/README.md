# gRPC Usage

### 目标
基于grpc的双向通信demo

### protoc (Protocol buffer compiler) 安装

Install `protoc` from [protobuf](https://github.com/protocolbuffers/protobuf/releases)
```shell
protoc --version  # Ensure compiler version is 3+
```

Install `Go plugins` for the `protoc`

```shell
# Install the protocol compiler plugins for Go using the following commands:
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

# Update your PATH so that the protoc compiler can find the plugins:
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 构建 gRPC 服务

定义 `.proto` 文件，本文以 `pixiu.proto` 为例

```protobuf
syntax="proto3";

option go_package = "go-learning/practise/grpc-practise/tunnel/tunnel";

package tunnel;

service Tunnel {
  // Client 调用此方法建立连接
  rpc Connect(stream Request) returns (stream Response);
}

message Request {
  string type = 1;  // "client_call" 或 "server_call"
  bytes payload = 2;
}

message Response {
  bytes result = 1;
}
```

### 生成 `gRPC` 代码

```shell
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  tunnel/tunnel.proto
```

执行命令之后，会在 `tunnel` 目录生成 `tunnel.pb.go` 和 `tunnel_grpc.pb.go` 代码文件
- `tunnel.pb.go` 结构体
- `tunnel_grpc.pb.go`: 客户端和服务端代码

### 实现 gRPC 服务端
```
```

### 实现 gRPC 客户端
```
```

### 执行

启动 `gRPC server`
``` shell
go run server.go
```

执行 `gRPC client`
``` shell
go run client.go

# 回显
2022/03/20 19:43:13 say hello caoyingjun 12345
```
