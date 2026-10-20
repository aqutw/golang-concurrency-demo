// Intentionally leaks one sender until process exit.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

//go:noinline
func leakySender(ch chan<- string, release <-chan struct{}) { <-release; ch <- "結果資料" }
func main() {
	ch := make(chan string)
	release := make(chan struct{})
	go leakySender(ch, release)
	select {
	case result := <-ch:
		fmt.Println(result)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timeout")
	}
	close(release)
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		buf := make([]byte, 64<<10)
		n := runtime.Stack(buf, true)
		stack := string(buf[:n])
		for _, frame := range strings.Split(stack, "\n\n") {
			if strings.Contains(frame, "main.leakySender(") && strings.Contains(frame, "[chan send]") {
				fmt.Println("LEAK: blocked sender")
				fmt.Println(frame)
				return
			}
		}
		runtime.Gosched()
	}
	fmt.Fprintln(os.Stderr, "sender did not reach blocked state")
	os.Exit(1)
}
