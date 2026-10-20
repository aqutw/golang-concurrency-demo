package main

import (
	"fmt"
	"sync"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false
	waiting := make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cond.L.Lock()
			defer cond.L.Unlock()
			waiting <- struct{}{}
			for !ready {
				cond.Wait()
			} // 喚醒後仍須持鎖重查共享條件。
			fmt.Printf("訂閱者 %d 收到通知，開始執行任務\n", id)
		}(i)
	}
	for i := 0; i < 3; i++ {
		<-waiting
	}
	cond.L.Lock()
	ready = true
	fmt.Println("廣播：條件已達成！")
	cond.Broadcast()
	cond.L.Unlock()
	wg.Wait()
}
