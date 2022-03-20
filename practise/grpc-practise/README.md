# gRPC Usage

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

option go_package = "go-learning/practise/grpc-practise/pixiu/pixiu";

package pixiu;

service Pixiu {
  rpc GetPixiu (PixiuRequest) returns (PixiuReply) {}
  // 其他接口
}

// The request message.
message PixiuRequest {
  int64  id = 1;
  string name = 2;
}

// The response message.
message PixiuReply {
  string message = 1;
}
```

### 生成 `gRPC` 代码

```shell
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  pixiu/pixiu.proto
```

执行命令之后，会在 `pixiu` 目录生成 `pixiu.pb.go` 和 `pixiu_grpc.pb.go` 代码文件
- `pixiu.pb.go` 结构体
- `pixiu_grpc.pb.go`: 客户端和服务端代码

### 实现 gRPC 服务端
``` go
...
type server struct {
	pd.UnimplementedPixiuServer
}

func (s *server) GetPixiu(ctx context.Context, in *pd.PixiuRequest) (*pd.PixiuReply, error) {
	log.Printf("Received %s %d", in.Name, in.Id)
	return &pd.PixiuReply{Message: fmt.Sprintf("%s %d", in.GetName(), in.GetId())}, nil
}

func main() {
	l, _ := net.Listen("tcp", ":30000")

	s := grpc.NewServer()
	pd.RegisterPixiuServer(s, &server{})

	if err = s.Serve(l); err != nil {
		...
	}
}
```

### 实现 gRPC 客户端
``` go
func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect rpc server %v", err)
	}
	defer conn.Close()

	c := pd.NewPixiuClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetPixiu(ctx, &pd.PixiuRequest{Id: 12345, Name: "caoyingjun"})
	...
}
```

### 执行

启动 `gRPC server`
``` shell
go run server/server.go
```

执行 `gRPC client`
``` shell
go run client/client.go

# 回显
2022/03/20 19:43:13 say hello caoyingjun 12345
```
