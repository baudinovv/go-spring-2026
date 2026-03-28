package main

import (
	"fmt"
	"sync"
)

// Solution 1: Using sync.Map 

func solutionSyncMap() {
	fmt.Println("=== Solution 1: sync.Map ===")

	var sm sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			sm.Store("key", key)
		}(i)
	}

	wg.Wait()

	value, ok := sm.Load("key")
	if ok {
		fmt.Printf("Value: %v\n", value)
	}
}

//  Solution 2: sync.RWMutex with a regular map 

func solutionRWMutex() {
	fmt.Println("=== Solution 2: sync.RWMutex ===")

	unsafeMap := make(map[string]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			mu.Lock()
			unsafeMap["key"] = key
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	mu.RLock()
	value := unsafeMap["key"]
	mu.RUnlock()

	fmt.Printf("Value: %d\n", value)
}

func main() {
	solutionSyncMap()
	solutionRWMutex()
}
