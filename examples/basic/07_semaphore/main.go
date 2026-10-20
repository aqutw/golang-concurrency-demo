package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 限制同時最多隻有 2 個 Goroutine 能執行
	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup
	var active atomic.Int32

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}
			if n := active.Add(1); n > 2 {
				panic("semaphore limit exceeded")
			} // 獲取令牌（若滿了則阻塞）

			fmt.Printf("任務 %d 進入臨界區\n", id)
			time.Sleep(20 * time.Millisecond) // 模擬執行

			active.Add(-1)
			<-sem // 釋放令牌
		}(i)
	}
	wg.Wait()
	fmt.Println("所有任務完成，併發上限=2")
}
