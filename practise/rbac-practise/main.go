package main

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 要使用自己定义的数据库rbac_db,最后的true很重要.默认为false,使用缺省的数据库名casbin,不存在则创建
	ada := gormadapter.NewAdapter("mysql", "root:password123456@tcp(pixiu01:3306)/casbin", true)
	enforcer := casbin.NewEnforcer("./model.conf", ada)

	// 从DB加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
	r := gin.Default()

	r.POST("/api/v1", func(c *gin.Context) {
		enforcer.AddPolicy("admin", "/api/v1/test", "GET")
		fmt.Println("Add policy")
	})

	r.DELETE("/api/v1", func(c *gin.Context) {
		enforcer.RemovePolicy("admin", "/api/v1/test", "GET")
		fmt.Println("delete Policy")
	})

	r.GET("/api/v1", func(c *gin.Context) {
		policies := enforcer.GetPolicy()
		for _, p := range policies {
			fmt.Println("policy", p)
		}
	})

	//使用自定义拦截器中间件
	r.Use(Authorization(enforcer))
	//创建请求
	r.GET("/api/v1/test", func(c *gin.Context) {
		fmt.Println("get v1 请求通过")
	})
	r.POST("/api/v1/test", func(c *gin.Context) {
		fmt.Println("post v1 请求通过")
	})

	_ = r.Run()
}

func Authorization(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		sub := "admin"

		if e.Enforce(sub, obj, act) {
			fmt.Println("验证通过")
			c.Next()
		} else {
			fmt.Println("无权访问")
			c.Abort()
			c.String(http.StatusUnauthorized, "无权访问")
		}
	}
}
