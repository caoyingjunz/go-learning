package main

import (
	"errors"
	"fmt"
	"reflect"
)

// 自定义 error
// 1. 使用 errors 中的new，直接return errors.New("erro msg")
// 2. fmt.Errorf()
// 3. 使用结构体去实现

var err1 error = errors.New("error msg1")
var err2 error = fmt.Errorf("error msg2")

func test1() error {
	return err1
}

func test2() error {
	return err2
}

type ErrStruct struct {
	Code int
	Msg  string
}

func (e *ErrStruct) Error() string {
	return e.Msg
}

func New(code int, mgs string) *ErrStruct {
	return &ErrStruct{
		Code: code,
		Msg:  mgs,
	}
}

func Test4() error {
	return New(444, "not fould test4")

}

func main() {

	err1 := test1()
	if err1 != nil {
		fmt.Println(err1.Error())
		fmt.Println(reflect.TypeOf(err1))
	}

	err2 := test2()
	if err2 != nil {
		fmt.Println(err2.Error())
		fmt.Println(reflect.TypeOf(err2))
	}

	err3 := New(404, "not fould test3")
	if err3 != nil {
		fmt.Println(err3.Error())
		fmt.Println(err3.Code)
		fmt.Println(reflect.TypeOf(err3))
	}

	err4 := Test4()
	if err4 != nil {
		fmt.Println(err4.Error())
	}
}
