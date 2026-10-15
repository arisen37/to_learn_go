package lab

import (
	"context"
	"errors"
	"sync"
)

func Map(ctx context.Context, nums []int, workers int, fn func(context.Context, int) (int, error)) ([]int, error) {
	if workers <= 0 || fn == nil {
		return nil, errors.New("invalid workers or callback")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	out := make([]int, len(nums))
	jobs := make(chan int)
	var wg sync.WaitGroup
	var once sync.Once
	var first error
	fail := func(err error) { once.Do(func() { first = err; cancel() }) }
	if workers > len(nums) {
		workers = len(nums)
	}
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case i, ok := <-jobs:
					if !ok {
						return
					}
					if ctx.Err() != nil {
						return
					}
					v, err := fn(ctx, nums[i])
					if err != nil {
						fail(err)
						return
					}
					out[i] = v
				}
			}
		}()
	}
send:
	for i := range nums {
		select {
		case <-ctx.Done():
			break send
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()
	if first != nil {
		return nil, first
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
