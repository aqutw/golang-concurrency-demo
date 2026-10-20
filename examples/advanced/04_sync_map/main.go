package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map
	var wg sync.WaitGroup

	// 併發寫入
	keys := []string{"apple", "banana", "cherry"}
	for _, k := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			sm.Store(key, len(key)) // 寫入鍵值對
		}(k)
	}
	wg.Wait()

	// 讀取
	val, ok := sm.Load("banana")
	if ok {
		fmt.Println("找到 banana, 長度為:", val)
	}

	// 走訪 Map (Range)
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // 返回 true 繼續迭代，返回 false 停止
	})
}
