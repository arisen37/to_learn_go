package main

import (
	"context"
	"fmt"
	"sync"
)

func run(ctx context.Context) (int, error) {
	jobs := make(chan int, 2)
	results := make(chan int)
	var workers sync.WaitGroup
	var producer sync.WaitGroup
	producer.Add(1)
	go func() {
		defer producer.Done()
		defer close(jobs)
		for n := 1; n <= 5; n++ {
			select {
			case jobs <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	for i := 0; i < 3; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for n := range jobs {
				select {
				case results <- n * n:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	joined := make(chan struct{})
	go func() { workers.Wait(); producer.Wait(); close(results); close(joined) }()
	total := 0
	for n := range results {
		total += n
	}
	<-joined
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return total, nil
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sum, err := run(ctx)
	fmt.Println("sum of squares:", sum, "error:", err)
}
