package lab

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	var active, peak atomic.Int64
	in := []int{1, 2, 3, 4, 5, 6}
	fn := func(ctx context.Context, n int) (int, error) {
		v := active.Add(1)
		defer active.Add(-1)
		for {
			old := peak.Load()
			if v <= old || peak.CompareAndSwap(old, v) {
				break
			}
		}
		time.Sleep(time.Millisecond)
		return n * n, nil
	}
	got, err := Map(context.Background(), in, 2, fn)
	if err != nil || !slices.Equal(got, []int{1, 4, 9, 16, 25, 36}) {
		t.Fatal(got, err)
	}
	if peak.Load() > 2 || active.Load() != 0 {
		t.Fatal("unbounded or unjoined workers")
	}
	if !slices.Equal(in, []int{1, 2, 3, 4, 5, 6}) {
		t.Fatal("mutated input")
	}
	boom := errors.New("boom")
	if got, err = Map(context.Background(), in, 3, func(context.Context, int) (int, error) { return 0, boom }); got != nil || !errors.Is(err, boom) {
		t.Fatal(got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = Map(ctx, in, 2, fn); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	for _, w := range []int{0, -1} {
		if _, err = Map(context.Background(), in, w, fn); err == nil {
			t.Fatal("invalid workers")
		}
	}
	if _, err = Map(context.Background(), in, 1, nil); err == nil {
		t.Fatal("nil callback")
	}
	if got, err = Map(context.Background(), nil, 2, fn); err != nil || len(got) != 0 {
		t.Fatal(got, err)
	}
}
func TestMapMidFlightCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := Map(ctx, []int{1, 2, 3}, 1, func(ctx context.Context, n int) (int, error) { close(started); <-ctx.Done(); return 0, ctx.Err() })
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("worker never started")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("cancellation did not join workers")
	}
}
