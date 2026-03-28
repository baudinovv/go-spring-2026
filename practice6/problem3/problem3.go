package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func startServer(ctx context.Context, name string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
				out <- fmt.Sprintf("[%s] metric: %d", name, rand.Intn(100))
			}
		}
	}()
	return out
}

// FanIn merges multiple input channels into a single output channel.
// It closes the result channel only after all input channels are closed.
// Supports cancellation via context.Context.
func FanIn(ctx context.Context, channels ...<-chan string) <-chan string {
	result := make(chan string)
	var wg sync.WaitGroup

	// Forward all messages from a single channel into result
	forward := func(ch <-chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}
				select {
				case result <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go forward(ch)
	}

	// Close result channel once all input goroutines finish
	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := startServer(ctx, "Alpha")
	ch2 := startServer(ctx, "Beta")
	ch3 := startServer(ctx, "Gamma")

	ch4 := FanIn(ctx, ch1, ch2, ch3)

	for val := range ch4 {
		fmt.Println(val)
	}
}
