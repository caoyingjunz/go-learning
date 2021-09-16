package main

import (
	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/endpoint"
	"go-learning/practise/gin-practise/middleware"
)

//GET http://127.0.0.1:8000/practise/get
//
//POST http://127.0.0.1:8000/practise/post
//
//{
//"name": "caoyingjun",
//"obj": {"k1": "v1", "k2": {"subk2": "subv2"}, "k3": "v3"},
//"pers": {"id": 123, "age": 19, "sex": "boy", "max": {"k1": "v1"}, "lis": ["1", "2", "3"]}
//}

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(middleware.LoggerToFile(), middleware.Auth)

	p := r.Group("/practise")
	{
		p.GET("/get", endpoint.GetPractise)
		p.POST("/post", endpoint.PostPractise)
	}

	r.Run(":8000")
}
