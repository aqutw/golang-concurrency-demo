package main

import (
	"context"
	"fmt"
)

func worker(ctx context.Context, ch <-chan int, done chan<- struct{}) {
	defer close(done)
	select {
	case <-ctx.Done():
		fmt.Println("FIXED: receiver canceled")
	case value := <-ch:
		fmt.Println(value)
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{})
	go worker(ctx, make(chan int), done)
	cancel()
	<-done
}
