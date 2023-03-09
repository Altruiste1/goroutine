package main

import (
	"fmt"
	"time"
)

func main() {
	pool := New(10)
	for i := 1; i < 10000; i++ {
		pool.NewTask(func() {
			fmt.Println(i)
			time.Sleep(10 * time.Second)
		})
	}
	// 保证所有的协程都执行完毕
	time.Sleep(time.Hour)
}
func (p *Pool) NewTask(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
}
func (p *Pool) worker(task func()) {
	defer func() { <-p.sem }()
	for {
		task()
		task = <-p.work
	}
}

type Pool struct {
	work chan func()   // 任务
	sem  chan struct{} // 数量
}

func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}
