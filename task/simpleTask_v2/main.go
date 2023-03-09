package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type task func()

var maxTask = 100
var chTask = make(chan taskF, maxTask)

type taskF struct {
	f       task
	addTime time.Time //添加时间
	name    int
}

type TaskQueue struct {
	queue []*taskF
	lock  sync.Mutex
}

var CurQueue = new(TaskQueue)

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
	go consumer1()
	for i := 0; i < maxTask; i++ {
		go consumer()
	}
}
func producer() {
	for i := 0; i < 10000; i++ {
		go func(i int) {
			f := taskF{
				f: func() {
					fmt.Println(i)
				},
				name:    i,
				addTime: time.Now(),
			}
			CurQueue.Push(&f)
		}(i)
	}
}

func consumer1() {
	for {
		f := CurQueue.Pop()
		if f != nil {
			chTask <- *f
		}
	}
}
func check() {
	for {
		fmt.Println("cur goroutine", runtime.NumGoroutine())
		fmt.Println("cur task", len(CurQueue.queue))
		time.Sleep(time.Second * 5)
	}

}

func consumer() {
	for {
		select {
		case a := <-chTask:
			fmt.Println("do task: ", a.name)
			a.f()
			fmt.Println("task finish ", a.name)
			time.Sleep(time.Second * 5)
		default:
			fmt.Println("wait 2 second")
			time.Sleep(time.Second * 10)
		}
	}
}

//加入元素
//这里需要使用指针
func (q *TaskQueue) Push(v *taskF) {
	q.lock.Lock()
	q.queue = append(q.queue, v)
	q.lock.Unlock()
}

//移除元素
func (q *TaskQueue) Pop() *taskF {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 {
		return nil
	}
	head := q.queue[0]
	if len(q.queue) > 1 {
		q.queue = q.queue[1:]
	} else {
		q.queue = nil
	}

	return head
}
