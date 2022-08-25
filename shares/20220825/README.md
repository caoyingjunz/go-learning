# Gin Web 框架的使用

## Gin 是什么
[Gin](https://github.com/gin-gonic/gin) is a web framework written in Go (Golang)

## Installation
```shell
go get -u github.com/gin-gonic/gin
```

## 最简单的 demo

```shell
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	_ = r.Run()
}
```
[demo](./demo.go)

## Gin 的参数处理

- 从 URL 中获取参数 (GET DELETE)
  - query 参数
    ```shell
    GET http://127.0.0.1:8080/v1/detail?name=caoyingjun&age=18
    ```
    从 c.Query 中获取参数
  
    ``` shell
    name := c.Query("name")
    age := c.Query("age")
    // do somethings
    ```

  - path参数
    ```shell
    GET http://127.0.0.1:8080/v1/name/caoyingjun/age/18
    ```
    从 c.Param 中获取参数
    ```shell
    name := c.Param("name")
    // do somethings
    ```

  - 从 body 中获取参数 (POST PUT)
  
    - 通常是构造和参数匹配的结构体，然后进行 bind
  
    ```shell
    p := struct {
          Name string `json:"name,omitempty"`
          Age  int    `json:"age,omitempty"`
    }{}
    _ = c.ShouldBindJSON(&p)
    // do somethings
    ```
[demo](./parameter.go)
  
## Gin 的返回值处理
构造期望访问的结构体

```shell
type Response struct {
	Code    int         `json:"code"` // 业务 code
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
}
```

然后根据请求的处理结果进行赋值即可
```shell
r := httputils.NewResponse()
p := struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}{}

_ = c.ShouldBindJSON(&p)
r.Result = p

httputils.SetSuccess(c, r)
```

## Gin 的中间件
通过在完成 gin 的初始化后，使用 Use 方法

```shell
r := gin.Default()

r.Use(LoggerToFile(), Auth) // 中间件
```
[demo](./parameter.go)

## 更多使用场景
[more](https://github.com/gin-gonic/gin#gin-web-framework)