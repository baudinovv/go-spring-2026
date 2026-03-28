package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//  Solution 1: sync.Mutex 

func solutionMutex() {
	fmt.Println("=== Solution 1: sync.Mutex ===")

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", counter)
}

//  Solution 2: sync/atomic 

func solutionAtomic() {
	fmt.Println("=== Solution 2: sync/atomic ===")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("Counter:", atomic.LoadInt64(&counter))
}

func main() {
	// The goroutines perform a non-atomic read-modify-write on `counter`,
	// so concurrent goroutines can read the same stale value and overwrite
	// each other's increments, causing lost updates — a classic data race.

	solutionMutex()
	solutionAtomic()
}
