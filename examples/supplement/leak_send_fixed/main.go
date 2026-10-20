package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)
	release, done := make(chan struct{}), make(chan struct{})
	go func() { defer close(done); <-release; ch <- "結果資料" }()
	select {
	case result := <-ch:
		fmt.Println(result)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timeout")
	}
	close(release)
	<-done
	fmt.Println("FIXED: buffered sender exited")
}
