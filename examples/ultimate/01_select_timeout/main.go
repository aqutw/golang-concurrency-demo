package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)
	done := make(chan struct{})
	release := make(chan struct{})
	go func() { defer close(done); <-release; ch <- "API 數據獲取成功" }()
	select {
	case res := <-ch:
		fmt.Println(res)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("錯誤: API 請求超時，中斷操作！")
	}
	close(release)
	<-done
	fmt.Println("緩衝通道讓逾時後的發送端仍可結束")
}
