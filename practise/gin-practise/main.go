package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-learning/practise/gin-practise/endpoint"
	"go-learning/practise/gin-practise/middleware"
)

var https = `
GET http://127.0.0.1:8000/practise/get

POST http://127.0.0.1:8000/practise/post

{
"name": "caoyingjun",
"obj": {"k1": "v1", "k2": {"subk2": "subv2"}, "k3": "v3"},
"pers": {"id": 123, "age": 19, "sex": "boy", "max": {"k1": "v1"}, "lis": ["1", "2", "3"]}
}

POST http://127.0.0.1:8000/practise/queue?queue=test

POST http://127.0.0.1:8000/practise/queue/after?after=after

GET http://127.0.0.1:8000/practise/pod?name=mariadb-0&namespace=kubez-sysns&key=config
`

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(middleware.LoggerToFile(), middleware.Auth)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go endpoint.WorkerSet.Run(2, stopCh)
	go endpoint.KubeEngine.Start(stopCh)

	p := r.Group("/practise")
	{
		p.GET("/get", endpoint.GetPractise)
		p.POST("/post", endpoint.PostPractise)
		p.POST("/queue", endpoint.TestQueue)
		p.POST("/queue/after", endpoint.TestAfterQueue)
		p.GET("/pod", endpoint.TestPod)
		p.POST("/optimise", endpoint.TestOptimise)
	}

	http.HandleFunc("/", downloadFile) //   设置访问路由
	go http.ListenAndServe(":8080", nil)

	_ = r.Run(":8000")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := r.Form["filename"]
	fmt.Println(fileName)
	path := "/Users/xxx/shuku.zip"
	http.ServeFile(w, r, path)
}
