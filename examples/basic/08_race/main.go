// ❌ 錯誤範例：會發生 Race Condition
package main

import (
	"fmt"
	"sync"
)

func main() {
	var count = 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count++ // 非安全併發寫入，結果通常不會是 1000
		}()
	}
	wg.Wait()
	fmt.Println("Count:", count)
}
