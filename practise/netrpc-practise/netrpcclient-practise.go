package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	var args = struct {
		A, B float32
	}{17, 8}

	var result = struct {
		Value float32
	}{}
	var client, err = rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("连接RPC服务失败：", err)
	}
	err = client.Call("MathService.Add", args, &result)
	if err != nil {
		fmt.Println("MathService.Add：", err)
	}
	fmt.Println("MathService.Add：", result.Value)

	err = client.Call("MathService.Divide", args, &result)
	if err != nil {
		fmt.Println("MathService.Divide：", err)
	}
	fmt.Println("MathService.Divide：", result.Value)
}
