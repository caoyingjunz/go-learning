package main

import "fmt"

// https://www.jianshu.com/p/b38b1719636e
// 如果说goroutine和channel是Go并发的两大基石，那么接口是Go语言编程中数据类型的关键。
// 在Go语言的实际编程中，几乎所有的数据结构都围绕接口展开，接口是Go语言中所有数据结构的核心。
// TODO: 对比context的源码，学习接口的正确使用方式

type DataWriter interface {
	WriteData(data interface{}) error
	CanWrite() bool
}

type dataWriter struct{}

func (d *dataWriter) WriteData(data interface{}) error {
	return nil
}

func (d *dataWriter) CanWrite() bool {
	return false
}

// 用于检查 dataWriter 是否实现了 DataWriter 接口
var _ DataWriter = &dataWriter{}

type file struct{}

func (f *file) WriteData(data interface{}) error {
	fmt.Println("writedata:", data)
	return nil
}

func (f *file) CanWrite() bool {
	return true
}

type Service interface {
	Start()
	Log(string)
}

type Logger struct {
}

func (l *Logger) Log(lg string) {
	fmt.Println(lg)
}

type GameService struct {
	Logger
}

func (g *GameService) Start() {
	fmt.Println("Starting now")
}

func main() {

	//f := new(file)
	//
	//var writer DataWriter
	//
	//writer = f
	//if writer.CanWrite() {
	//	writer.WriteData("data")
	//}

	// TODO, 暂时还没有想通 接口的优势
	var s Service = new(GameService)
	s.Log("1234")
	s.Start()
}
