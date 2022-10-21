package main

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	user := "root"
	password := "password123456"
	ip := "pixiu01"
	port := 3306
	database := "rbacs"
	dbConnection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", user, password, ip, port, database)
	// must declare the err to aviod panic: runtime error: invalid memory address or nil pointer dereferences
	var err error
	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	enforcer, err := casbin.NewEnforcer("./model.conf", adapter)
	if err != nil {
		panic(err)
	}

	// 从DB加载策略
	if err = enforcer.LoadPolicy(); err != nil {
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

		if ok, err := e.Enforce(sub, obj, act); err == nil && ok {
			fmt.Println("验证通过")
			c.Next()
		} else {
			fmt.Println("无权访问")
			c.Abort()
			c.String(http.StatusUnauthorized, "无权访问")
		}
	}
}
