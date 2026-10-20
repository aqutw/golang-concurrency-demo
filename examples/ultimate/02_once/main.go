package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Theme string
}

var (
	instance *Config
	once     sync.Once
	wg       sync.WaitGroup
)

func GetConfigInstance() *Config {
	// 確保無論多少個 Goroutine 同時呼叫，這段邏輯都只會執行一次
	once.Do(func() {
		fmt.Println("--> 正在讀取設定檔並初始化 Config 物件... (全域僅限一次)")
		instance = &Config{Theme: "Dark Mode"}
	})
	return instance
}

func main() {
	// 模擬 5個併發請求同時獲取設定檔
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := GetConfigInstance()
			fmt.Printf("Worker %d 獲取的配置: %s\n", id, cfg.Theme)
		}(i)
	}
	wg.Wait()
}
