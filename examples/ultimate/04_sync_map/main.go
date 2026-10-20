package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	var wg sync.WaitGroup

	// 1. 併發寫入 Store
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			m.Store(key, key*10) // 儲存鍵值對
		}(i)
	}
	wg.Wait()

	// 2. 併發讀取 Load
	wg.Add(1)
	go func() {
		defer wg.Done()
		if val, ok := m.Load(2); ok {
			fmt.Println("成功讀取 Key 2 的值:", val)
		}
	}()

	// 3. 遍歷 Range
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("鍵: %v, 值: %v\n", key, value)
		return true // 返回 true 繼續遍歷，返回 false 則停止
	})

	wg.Wait()
}
