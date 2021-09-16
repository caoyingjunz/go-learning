package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// make(T) 返回一个类型为 T 的初始值，它只适用于3种内建的引用类型：slice、map 和 channel
// https://blog.csdn.net/benben_2015/article/details/78917374

type Student struct {
	Name  interface{} `json:"name"`
	Age   interface{} `json:"age"`
	Sex   interface{}
	Class interface{} `json:"class"`
}

type Class struct {
	Name  string
	Grade int
}

func main() {

	var stu = Student{
		Name: "jun",
		Age:  18,
		Sex:  "Boy",
	}
	class := new(Class)
	class.Name = "1 class"
	class.Grade = 1
	stu.Class = class

	jsonStu, err := json.Marshal(stu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonStu))

	// json str 转 map
	jsonStr := `{"name":"jun","age":18,"Sex":"Boy","class":{"Name":"1 class","Grade":1}}`
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonData)
	fmt.Println(fmt.Sprintf("%T", jsonData))
	fmt.Println(jsonStr)

	// json str to struct
	var stud Student
	json.Unmarshal([]byte(jsonStr), &stud)
	fmt.Println(stud.Name)

	//array 到 json str
	arr := []string{"test1", "test2", "test3"}
	lang, err := json.Marshal(arr)
	if err == nil {
		fmt.Println(string(lang))
		fmt.Println(fmt.Sprintf("%T", string(lang)))
	}

	// json 到 []string
	var tt []string
	json.Unmarshal(lang, &tt)
	fmt.Println(tt)
}
