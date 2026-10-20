package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Data string
}

// worker 接收一個總控制的 ctx
func worker(id int, ctx context.Context, tasks <-chan Task) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("⚠️ Worker %d: 收到系統停止/超時訊號，結束退出。\n", id)
			return
		case task, ok := <-tasks:
			if !ok {
				fmt.Printf("Worker %d: 任務通道已關閉，正常結束。\n", id)
				return
			}

			// 針對「單一任務」設定超時保護（例如：每個任務最多執行 400 毫秒）
			taskCtx, cancel := context.WithTimeout(ctx, 400*time.Millisecond)

			fmt.Printf("Worker %d: 開始處理任務 %d\n", id, task.ID)
			err := processTask(taskCtx, task)
			cancel() // 釋放 context 資源

			if err != nil {
				fmt.Printf("❌ Worker %d: 任務 %d 失敗: %v\n", id, task.ID, err)
			} else {
				fmt.Printf("✅ Worker %d: 任務 %d 成功完成\n", id, task.ID)
			}
		}
	}
}

func processTask(ctx context.Context, t Task) error {
	// 模擬不同長度的耗時任務
	simulatedTime := 200 * time.Millisecond
	if t.ID == 2 {
		simulatedTime = 800 * time.Millisecond // 任務 2 故意讓他超時
	}

	timer := time.NewTimer(simulatedTime)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err() // 如果超時或被取消，直接中斷返回
	case <-timer.C:
		return nil // 順利執行完畢
	}
}

func main() {
	// 建立一個全域 2 秒超時的 Context 作為雙重保險
	rootCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tasks := make(chan Task, 5)

	var wg sync.WaitGroup
	// 啟動 2 個 Worker
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) { defer wg.Done(); worker(id, rootCtx, tasks) }(w)
	}

	// 派發 3 個任務
	for i := 1; i <= 3; i++ {
		tasks <- Task{ID: i, Data: "處理數據"}
	}
	close(tasks)

	wg.Wait()
	fmt.Println("所有 Worker 已退出")
}
