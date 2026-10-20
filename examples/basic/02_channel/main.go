package main

import "fmt"

func main() {
	// 建立一個無緩衝的 int Channel
	ch := make(chan int)

	go func() {
		ch <- 42 // 傳送資料到 Channel
	}()

	val := <-ch // 從 Channel 接收資料（會阻塞直到收到資料）
	fmt.Println("接收到的值:", val)
}
