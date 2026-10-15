package lab

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
)

type result struct {
	key   string
	value int64
	err   error
}

func Aggregate(parent context.Context, r io.Reader, workers int) (map[string]int64, error) {
	if workers <= 0 {
		return nil, errors.New("invalid workers")
	}
	if err := parent.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	rows := make(chan []string, workers)
	results := make(chan result)
	var wg sync.WaitGroup
	emit := func(v result) bool {
		select {
		case results <- v:
			return true
		case <-ctx.Done():
			return false
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(rows)
		c := csv.NewReader(r)
		head, err := c.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			emit(result{err: err})
			return
		}
		if len(head) != 2 || head[0] != "key" || head[1] != "value" {
			emit(result{err: errors.New("invalid header")})
			return
		}
		for {
			if ctx.Err() != nil {
				return
			}
			row, err := c.Read()
			if err == io.EOF {
				return
			}
			if err != nil {
				emit(result{err: err})
				return
			}
			select {
			case rows <- row:
			case <-ctx.Done():
				return
			}
		}
	}()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case row, ok := <-rows:
					if !ok {
						return
					}
					key := strings.TrimSpace(row[0])
					n, err := strconv.ParseInt(row[1], 10, 64)
					if err == nil && key == "" {
						err = errors.New("blank key")
					}
					if !emit(result{key, n, err}) {
						return
					}
					if err != nil {
						return
					}
				}
			}
		}()
	}
	joined := make(chan struct{})
	go func() { wg.Wait(); close(results); close(joined) }()
	out := make(map[string]int64)
	var first error
	for v := range results {
		if first != nil {
			continue
		}
		if v.err != nil {
			first = v.err
			cancel()
			continue
		}
		old := out[v.key]
		if v.value > 0 && old > math.MaxInt64-v.value || v.value < 0 && old < math.MinInt64-v.value {
			first = errors.New("sum overflow")
			cancel()
			continue
		}
		out[v.key] = old + v.value
	}
	<-joined
	if first != nil {
		return nil, first
	}
	if err := parent.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
