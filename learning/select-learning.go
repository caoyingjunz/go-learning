package main

import (
	"fmt"
	"time"
)

// 参考 https://www.jianshu.com/p/2a1146dc42c3

// 一个select语句用来选择哪个case中的发送或接收操作可以被立即执行。类似于switch语句，但是它的case涉及到channel有关的I/O操作。
// select就是用来监听和channel有关的IO操作，当 IO 操作发生时，触发相应的动作。

func main() {
	// 多个可操作时，随机选一个执行，不存在则阻塞
	//所有channel表达式都会被求值、所有被发送的表达式都会被求值。求值顺序：自上而下、从左到右.
	// break关键字结束select
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 3
	ch2 <- 5

	// 基本使用
	select {
	case a := <-ch1:
		fmt.Println("ch1 selected.", a)
		if a == 3 {
			break
		}
		fmt.Println("ch1 selected after break")
	case <-ch2:
		fmt.Println("ch2 selected.")
		fmt.Println("ch2 selected without break")
	}
	fmt.Println("OVER")

	// 超时实现
	// TODO 后续 context 实现超时

	// 官方用法
	//select {
	//case m := <-c:
	//	handle(m)
	//case <-time.After(10 * time.Second):
	//	fmt.Println("timed out")
	//}}

	ch := make(chan int, 1)

	select {
	case num := <-ch: // 读超时
		//case ch <- 123  // 写超时
		fmt.Println("test action", num)
	case <-time.After(5 * time.Second):
		fmt.Println("超时啦")
	}

}
