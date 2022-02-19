package log

type Configuration struct {
	LogFile  string
	LogLevel string

	RotateMaxSize    int
	RotateMaxAge     int
	RotateMaxBackups int
	Compress         bool
}

type Logger interface {
	Info(args ...interface{})
	Infof(f string, args ...interface{})
}

var Glog Logger
var AccessLog Logger

// 注册 AccessLog
func Register(logDir string) {
	afileLog := "/Users/caoyuan/workstation/go-learning/practise/gin-practise/practise-access.log"
	AccessLog, _ = NewZapLogger(Configuration{
		LogFile:          afileLog,
		LogLevel:         "INFO",
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})
}

// 或者 import 的时候初始化
func init() {
	fileLog := "/Users/caoyuan/workstation/go-learning/practise/gin-practise/practise-log.log"
	Glog, _ = NewZapLogger(Configuration{
		LogFile:          fileLog,
		LogLevel:         "INFO",
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})
}
