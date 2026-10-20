package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done() // 執行完畢後將計數器 -1
	fmt.Printf("Worker %d 啟動\n", id)
	time.Sleep(time.Millisecond * 500)
	fmt.Printf("Worker %d 完成\n", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // 計數器 +1
		go worker(i, &wg)
	}

	wg.Wait() // 阻塞，直到計數器歸零
	fmt.Println("所有 Worker 執行完畢")
}
