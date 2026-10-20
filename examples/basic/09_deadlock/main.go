// Intentionally broken: run without -race to see Go's deadlock diagnostic.
package main

import "sync"

func main() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	aLocked, bLocked := make(chan struct{}), make(chan struct{})
	go func() {
		defer wg.Done()
		mu1.Lock()
		close(aLocked)
		<-bLocked
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()
	go func() {
		defer wg.Done()
		mu2.Lock()
		close(bLocked)
		<-aLocked
		mu1.Lock()
		mu1.Unlock()
		mu2.Unlock()
	}()
	wg.Wait() // 雙方都拿到第一把鎖後才索取第二把鎖，穩定重現死鎖。
}
