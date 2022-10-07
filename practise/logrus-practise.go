package main

// 参考文档，logrus 官方源码 https://github.com/sirupsen/logrus
// 打印行数和文件 https://github.com/sirupsen/logrus/issues/63
// 暂时够用，用空的话，可以测试下钩子，目前看起来用不到

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	//用日志实例的方式使用日志
	//log.Out = os.Stdout   //日志标准输出

	// 也可以把log 写入日志文件
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("打开日志文件失败")
	}

	//log.Out = file
	log.SetOutput(file)

	// Only log the warning severity or above.
	log.SetLevel(logrus.WarnLevel)
}

func main() {

	// 是否需要打印行数和文件名称（影响性能），也可以对比较长的文件名进行切割，去除不必要的目录
	//log.SetReportCaller(true)

	log.WithFields(logrus.Fields{"request_id": "666"}).Info("测试 info ")
	log.WithFields(logrus.Fields{"request_id": "888"}).Warn("测试 warn ")
	log.WithFields(logrus.Fields{"request_id": "888"}).Error("测试 error ")
}
