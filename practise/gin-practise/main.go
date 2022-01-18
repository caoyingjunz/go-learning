package main

import (
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
	//go endpoint.KubeEngine.Start(stopCh)

	p := r.Group("/practise")
	{
		p.GET("/get", endpoint.GetPractise)
		p.POST("/post", endpoint.PostPractise)
		p.POST("/queue", endpoint.TestQueue)
		p.POST("/queue/after", endpoint.TestAfterQueue)
		p.GET("/pod", endpoint.TestPod)
		p.POST("/optimise", endpoint.TestOptimise)
		p.GET("/download", Download)
	}

	http.HandleFunc("/download", DownloadFile)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	_ = r.Run(":8000")
}

func Download(c *gin.Context) {
	filename := c.Query("filename")

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename+".tar")
	c.File("/Users/xxx/yyy.tar")
}

func DownloadFile(w http.ResponseWriter, req *http.Request) {
	fileName := req.URL.Query().Get("filename")
	if len(fileName) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("get file name failed"))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+".tar")
	w.Header().Set("Content-Type", "application/octet-stream")

	path := "/Users/xxx/yyy.tar"
	http.ServeFile(w, req, path)
}
