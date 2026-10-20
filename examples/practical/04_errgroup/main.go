package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	firstDone := make(chan struct{})
	g.Go(func() error {
		defer close(firstDone)
		fmt.Println("任務 1：成功結束")
		return nil
	})
	g.Go(func() error {
		<-firstDone // 先展示成功，再觸發錯誤，避免依賴 Sleep 的排程假設。
		return errors.New("資料庫連線中斷")
	})
	g.Go(func() error {
		<-ctx.Done()
		fmt.Println("任務 3：已即時自我中斷並釋放資源")
		return ctx.Err()
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("核心錯誤為: %v\n", err)
	}
}
