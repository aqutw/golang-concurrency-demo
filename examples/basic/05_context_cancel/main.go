package main

import (
	"context"
	"fmt"
)

func scanProgress(ctx context.Context, started, done chan struct{}) {
	defer close(done)
	fmt.Println("正在掃描中...")
	close(started)
	<-ctx.Done()
	fmt.Println("收到取消通知，停止掃描。原因:", ctx.Err())
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started, done := make(chan struct{}), make(chan struct{})
	go scanProgress(ctx, started, done)
	<-started
	cancel()
	<-done // 等待真正結束，不以 Sleep 猜測生命週期。
}
