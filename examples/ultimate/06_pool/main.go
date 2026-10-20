package main

import (
	"bytes"
	"fmt"
	"sync"
)

// 建立一個快取 bytes.Buffer 的池
var bufferPool = sync.Pool{
	// 如果池子裡沒有快取物件，定義如何新增一個
	New: func() interface{} {
		fmt.Println("--> 池中無可用物件，配置新記憶體空間")
		return new(bytes.Buffer)
	},
}

func main() {
	// 1. 從池中獲取物件 (Get)
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.WriteString("Data row 1")
	fmt.Println("使用 buffer 內容:", buf.String())

	// 2. 使用完畢，重設（Reset）物件狀態，並還回池中 (Put)
	buf.Reset()
	bufferPool.Put(buf)

	// 3. 再次獲取；Pool 可隨時丟棄物件，因此可能再次觸發 New。
	buf2 := bufferPool.Get().(*bytes.Buffer)
	buf2.WriteString("Data row 2")
	fmt.Println("複用 buffer 內容:", buf2.String())

	// 記得放回
	buf2.Reset()
	bufferPool.Put(buf2)
}
