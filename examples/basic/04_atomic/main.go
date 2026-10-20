package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var count int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&count, 1) // 安全地原子加 1
		}()
	}

	wg.Wait()
	fmt.Println("最終計數:", atomic.LoadInt64(&count))
}
