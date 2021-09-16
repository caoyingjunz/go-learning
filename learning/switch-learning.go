package main

import "fmt"

// Go语言改进了 switch 的语法设计，case 与 case 之间是独立的代码块，不需要通过 break 语句跳出当前 case 代码块以避免执行到下一行

func main() {

	var str = "test"
	switch str {
	case "123":
		fmt.Println("123")
	case "test":
		fmt.Println("test")
	default:
		fmt.Println("default")
	}

	// 也可以支持一分值多个值
	var str2 = "aaaa"
	switch str2 {
	case "123", "234":
		fmt.Println("多个值")
	default:
		fmt.Println("defualt")
	}

	// 支持分支表达式
	var num = 7
	switch { // switch后面 可以不加值
	case num > 8 && num < 5:
		fmt.Println("aaa")
	default:
		fmt.Println("bbb")
	}

	// 处理指针的用法
	var t interface{} // 重点，后续用的比较多
	t = 10
	switch tv := t.(type) {
	case int:
		fmt.Println("答应t的值", t, "类型为", fmt.Sprintf("%T", tv))
	case string:
		fmt.Println("SSSS")
	}
}
