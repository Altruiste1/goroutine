## v1
最简单的task
producer一开始生产10000个任务
根据设定的并发数创建数量为并发数的consumer
检查剩余的goroutine

## v2
在v1基础上,检测仅与task相关的goroutine数量
添加task的一些信息，添加时间，剩余task数等
