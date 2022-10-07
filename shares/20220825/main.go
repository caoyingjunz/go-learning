package main

import (
	"fmt"

	"net/http"

	"github.com/caoyingjunz/gopixiu/api/server/httputils"
	"github.com/gin-gonic/gin"
	_ "go-learning/shares/20220825/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// LoggerToFile 日志
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("log")
	}
}

func main() {
	r := gin.Default()

	r.Use(LoggerToFile(), Limiter, Auth) // 中间件
	g1 := r.Group("/v1")
	{
		g1.GET("/detail", getParametersFromQuery)
		g1.GET("/name/:name/age/:age", getParametersFromPath)
		g1.POST("/create", getParametersFromBody)
	}

	// TODO: the other groups
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	_ = r.Run()
}

// @BasePath /v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func getParametersFromQuery(c *gin.Context) {
	r := httputils.NewResponse()

	// do something

	r.Result = map[string]interface{}{"name": c.Query("name"), "age": c.Query("age")}
	httputils.SetSuccess(c, r)
}

func getParametersFromPath(c *gin.Context) {
	r := httputils.NewResponse()

	// do something

	r.Result = map[string]interface{}{"name": c.Param("name"), "age": c.Param("age")}
	httputils.SetSuccess(c, r)
}

func getParametersFromBody(c *gin.Context) {
	r := httputils.NewResponse()
	p := struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}{}
	_ = c.ShouldBindJSON(&p)

	// do something

	r.Result = p
	httputils.SetSuccess(c, r)
}
