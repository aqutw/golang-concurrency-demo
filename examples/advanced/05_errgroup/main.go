package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	ready := make(chan struct{}, 2)
	for _, id := range []int{1, 3} {
		id := id
		g.Go(func() error {
			ready <- struct{}{}
			<-ctx.Done()
			fmt.Printf("Task %d: 收到群組取消通知\n", id)
			return ctx.Err()
		})
	}
	g.Go(func() error { <-ready; <-ready; return errors.New("資料庫連線中斷") })
	if err := g.Wait(); err != nil {
		fmt.Println("最終群組錯誤:", err)
	}
}
