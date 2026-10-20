package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Server string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfigInstance() *Config {
	// 確保內部的匿名函式在整個程式生命週期中只被執行一次
	once.Do(func() {
		fmt.Println("--> 正在初始化設定檔（僅此一次）...")
		instance = &Config{Server: "localhost:8080"}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := GetConfigInstance()
			fmt.Printf("Worker %d 取得設定: %s\n", id, cfg.Server)
		}(i)
	}
	wg.Wait()
}
