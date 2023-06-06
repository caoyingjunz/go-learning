package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/caoyingjunz/pixiu/api/server/httputils"
)

// Auth 认证
func Auth(c *gin.Context) {
	fmt.Println("auth")
}

// Limiter 限速
func Limiter(c *gin.Context) {
	var limit bool
	if limit {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "message": "服务器繁忙，请稍后再试"})
		return
	}

	fmt.Println("limiter")
}

func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		}{}

		_ = ShouldBindWith(c, &obj)

		c.Next()
		fmt.Println(c.Value("name"))
	}
}

// LoggerToFile 日志
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("log")
	}
}

func ShouldBindWith(c *gin.Context, v interface{}) error {
	rawData, err := c.GetRawData()
	if err != nil {
		return err
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
	if err = json.Unmarshal(rawData, v); err != nil {
		return err
	}

	return nil
}

type Params struct {
	Name string `json:"name,omitempty" uri:"name" binding:"required" form:"name"`
	Age  int    `json:"age,omitempty" uri:"age" binding:"required" form:"age"`
}

type Objects struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name string `json:"name"`
}

func main() {
	r := gin.Default()

	r.Use(LoggerToFile(), Limiter, Auth, Audit()) // 中间件
	g1 := r.Group("/v1")
	{
		g1.GET("/detail", getParametersFromQuery)
		g1.GET("/name/:name/age/:age", getParametersFromPath)
		g1.POST("/create", getParametersFromBody)
		g1.PUT("/update", updateObject)
	}

	// TODO: the other groups
	_ = r.Run()
}

func updateObject(c *gin.Context) {
	r := httputils.NewResponse()
	var p Objects
	if err := c.ShouldBindJSON(&p); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	r.Result = p
	httputils.SetSuccess(c, r)
}

func getParametersFromQuery(c *gin.Context) {
	r := httputils.NewResponse()
	var p Params
	_ = c.ShouldBindQuery(&p)

	r.Result = map[string]interface{}{"name": c.Query("name"), "age": c.Query("age")}
	httputils.SetSuccess(c, r)
}

func getParametersFromPath(c *gin.Context) {
	r := httputils.NewResponse()

	// do something
	var p Params
	c.ShouldBindQuery(&p)
	// do something

	r.Result = map[string]interface{}{"name": c.Param("name"), "age": c.Param("age")}
	httputils.SetSuccess(c, r)
}

func getParametersFromBody(c *gin.Context) {
	r := httputils.NewResponse()
	var p Params
	if err := c.ShouldBindJSON(&p); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	// do something
	r.Result = p
	httputils.SetSuccess(c, r)

	c.Set("name", "caoyingjunz")
}
