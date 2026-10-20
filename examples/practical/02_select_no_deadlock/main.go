package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int) // 無緩衝通道

	// --- 情境 A：使用 default 進行「嘗試性」寫入，若失敗立馬走人 ---
	select {
	case ch <- 100:
		fmt.Println("成功寫入數據")
	default:
		// 因為此時沒有其他 Goroutine 準備好接收 ch，所以會直接走這裡，成功避免死鎖！
		fmt.Println("⚠️ 警告：目前通道無人接收，放棄寫入以避免死鎖")
	}

	// --- 情境 B：使用 time.After 進行「限時等待」接收 ---
	select {
	case val := <-ch:
		fmt.Println("接收到數據:", val)
	case <-time.After(500 * time.Millisecond):
		// 等了 500 毫秒還是沒有人送數據進來，自動中斷，成功避免死鎖！
		fmt.Println("⏰ 接收數據超時，安全退出")
	}
}
