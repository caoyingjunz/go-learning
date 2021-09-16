package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"go-learning/practise/gin-practise/log"
)

//refer to https://mp.weixin.qq.com/s/gBWEHe20Lv_2wBSlM2WeVA

func Auth(c *gin.Context) {
	// TODO: 用于进行 auth 校验
	test := c.Query("test")
	if test == "err" {
		c.AbortWithStatusJSON(400, map[string]string{"error": "log"})
		return
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
