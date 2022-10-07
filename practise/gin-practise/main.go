package main

import (
	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/endpoint"
	"go-learning/practise/gin-practise/log"
	"go-learning/practise/gin-practise/middleware"
	"go-learning/practise/gin-practise/worker"
)

var https = `
###
GET http://127.0.0.1:8000/practise/get

###
POST http://127.0.0.1:8000/practise/post

{
	"name": "caoyingjun",
	"obj": {"k1": "v1", "k2": {"subk2": "subv2"}, "k3": "v3"},
	"pers": {"id": 123, "age": 19, "sex": "boy", "max": {"k1": "v1"}, "lis": ["1", "2", "3"]}
}

###
POST http://127.0.0.1:8000/practise/queue?queue=test

###
POST http://127.0.0.1:8000/practise/queue/after?after=after

###
GET http://127.0.0.1:8000/practise/pod?name=mariadb-0&namespace=kubez-sysns&key=config
`

type options struct {
	addr   string
	engine *gin.Engine
	ws     worker.WorkerInterface
	logDir string
}

func (c *options) registerHttpRoute() {
	c.engine.Use(middleware.LoggerToFile(), middleware.Auth)

	m := c.engine.Group("middleware", middleware.AllowAccess())
	{
		m.POST("", endpoint.PostMid)
	}

	p := c.engine.Group("/practise")
	{
		p.GET("/get", endpoint.GetPractise)
		p.POST("/post", endpoint.PostPractise)
		p.POST("/queue", endpoint.TestQueue)
		p.POST("/queue/after", endpoint.TestAfterQueue)
		p.GET("/pod", endpoint.TestPod)
		p.POST("/optimise", endpoint.TestOptimise)
		p.GET("/download", endpoint.Download)
	}
}

func (c *options) run() {
	go func() {
		if err := c.engine.Run(c.addr); err != nil {
			panic(err)
		}
	}()

	go func() {
		stopCh := make(chan struct{})
		defer close(stopCh)
		c.ws.Run(2, stopCh)
	}()
}

func (c *options) registerLog() {
	log.Register(c.logDir)
}

func (c *options) registerController() {
	p := endpoint.New()
	endpoint.Register(p)
}

func (c *options) registerWorker() {
	c.ws = worker.NewWorker()
	endpoint.RegisterWorker(c.ws)
}

func NewHttpServer(addr string) *options {
	o := &options{
		addr:   addr,
		engine: gin.Default(),
	}

	o.registerWorker()

	o.registerLog()

	o.registerHttpRoute()
	return o
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	s := NewHttpServer(":8080")
	s.run()

	select {}
}
