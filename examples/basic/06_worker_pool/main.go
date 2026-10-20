package main

import (
	"fmt"
	"sort"
	"sync"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d 處理任務 %d\n", id, j)
		results <- j * 2
	}
}
func main() {
	jobs, results := make(chan int, 5), make(chan int, 5)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) { defer wg.Done(); worker(id, jobs, results) }(w)
	}
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	go func() { wg.Wait(); close(results) }()
	var values []int
	for value := range results {
		values = append(values, value)
	}
	sort.Ints(values)
	fmt.Println("結果:", values)
}
