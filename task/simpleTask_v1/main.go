package main

import (
	"fmt"
	"runtime"
	"time"
)

type task func()

var maxTask = 100
var chTask = make(chan task, maxTask)

//
func main() {
	// 生产10000个任务
	go producer()
	// 根据最大并发数创建消费端
	go mainConsumer()
	// 检查goroutine数量
	go check()
	time.Sleep(time.Hour)
}

func mainConsumer() {
	for i := 0; i < maxTask; i++ {
		go consumer()
	}
}
func producer() {
	for i := 0; i < 10000; i++ {
		go func(i int) {
			chTask <- func() {
				fmt.Println("cur", i)
			}
		}(i)
	}
}

func check() {
	for {
		fmt.Println("cur goroutine", runtime.NumGoroutine())
		time.Sleep(time.Second * 5)
	}

}

func consumer() {
	for {
		select {
		case a := <-chTask:
			a()
			time.Sleep(time.Second * 5)
		default:
			fmt.Println("wait 2 second")
			time.Sleep(time.Second * 10)
		}
	}
}
