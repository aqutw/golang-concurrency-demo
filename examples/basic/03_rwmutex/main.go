package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu    sync.RWMutex
	value map[string]int
}

func (c *SafeCounter) Inc(key string)       { c.mu.Lock(); defer c.mu.Unlock(); c.value[key]++ }
func (c *SafeCounter) Value(key string) int { c.mu.RLock(); defer c.mu.RUnlock(); return c.value[key] }
func main() {
	c := SafeCounter{value: make(map[string]int)}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Inc("somekey")
				_ = c.Value("somekey")
			}
		}()
	}
	wg.Wait()
	fmt.Println("RWMutex 計數:", c.Value("somekey"))
}
