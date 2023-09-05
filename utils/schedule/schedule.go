package schedule

import (
	"fmt"
	"time"
)

// 定时任务
func ScheduleTask(seconds int, task func()) {
	// 创建一个定时器，设置触发时间为指定秒数后
	timer := time.NewTimer(time.Duration(seconds) * time.Second)

	// 使用匿名 goroutine 等待定时器触发
	go func() {
		<-timer.C // 等待定时器的触发事件
		task()    // 执行传入的任务函数
	}()

	// 注意：如果不停止定时器，它将一直触发
	// 这里可以选择停止定时器，如果只需要执行一次任务
	timer.Stop()
}

// 暂停函数，接受秒数作为参数，并返回一个通道
func Sleep(seconds int) chan struct{} {
	// 创建一个通道用于控制暂停
	pauseCh := make(chan struct{})

	// 启动一个 goroutine 来执行定时暂停
	go func() {
		fmt.Printf("暂停 %d 秒\n", seconds)
		timer := time.NewTicker(1 * time.Second)
		defer timer.Stop()

		for i := seconds; i >= 1; i-- {
			select {
			case <-timer.C:
				fmt.Printf("%d ", i)
			case <-pauseCh:
				fmt.Println("\n提前结束暂停")
				return
			}
		}

		fmt.Println("\n暂停结束")
		close(pauseCh) // 结束暂停
	}()

	return pauseCh
}

// 结束暂停的函数，接受暂停通道作为参数
func EndPause(pauseCh chan struct{}) {
	close(pauseCh)
}
