package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/log"
)

const (
	testName = "test"
)

func Auth(c *gin.Context) {
	// TODO: 用于进行 auth 校验
	test := c.Query("test")
	if test == "err" {
		c.AbortWithStatusJSON(400, map[string]string{"error": "log"})
		return
	}
}

type requestObject struct {
	Name string `json:"name"`
}

func (o *requestObject) validate() error {
	return nil
}

func parseObjectFromRequest(c *gin.Context) (*requestObject, error) {
	data, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	var obj requestObject
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	if err = obj.validate(); err != nil {
		return nil, err
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 回写字节流
	return &obj, nil
}

// AllowAccess 从 gin.Context 获取指定字段
func AllowAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("allow access start")
		o, err := parseObjectFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(500, map[string]string{"error": err.Error()})
			return
		}

		if o.Name == testName {
			c.AbortWithStatusJSON(400, map[string]string{"error": "test name"})
			return
		}

		fmt.Println("access allow done")
	}
}

func LoggerToFile() gin.HandlerFunc {
	handlerFunc := func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求操作
		c.Next()

		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()

		log.AccessLog.Infof("| %3d | %13v | %15s | %s | %s |", statusCode, latencyTime, clientIp, reqMethod, reqUri)
	}
	return handlerFunc
}
