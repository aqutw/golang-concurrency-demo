package main

import (
	"fmt"
	"time"
)

func main() {
	ch1, ch2 := make(chan string), make(chan string)
	done := make(chan struct{})
	go func() {
		defer close(done)
		ch1 <- "來自 ch1 的資料"
		ch2 <- "來自 ch2 的資料"
	}()
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("接收成功:", msg)
		case msg := <-ch2:
			fmt.Println("接收成功:", msg)
		}
	}
	<-done
	// 發送端已結束；下一輪沒有資料，保證走超時分支。
	timer := time.NewTimer(20 * time.Millisecond)
	defer timer.Stop()
	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	case <-timer.C:
		fmt.Println("錯誤: 讀取超時！")
	}
}
