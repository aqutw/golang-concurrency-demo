// Intentionally leaks one receiver until process exit.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

//go:noinline
func leakyWorker(ch <-chan int) { fmt.Println(<-ch) }
func main() {
	ch := make(chan int)
	go leakyWorker(ch)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		buf := make([]byte, 64<<10)
		n := runtime.Stack(buf, true)
		for _, frame := range strings.Split(string(buf[:n]), "\n\n") {
			if strings.Contains(frame, "main.leakyWorker(") && strings.Contains(frame, "[chan receive]") {
				fmt.Println("LEAK: blocked receiver")
				fmt.Println(frame)
				return
			}
		}
		runtime.Gosched()
	}
	fmt.Fprintln(os.Stderr, "receiver did not reach blocked state")
	os.Exit(1)
}
