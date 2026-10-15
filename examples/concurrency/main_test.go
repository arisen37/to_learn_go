package main

import (
	"context"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	got, err := run(context.Background())
	if err != nil || got != 55 {
		t.Fatal(got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err = run(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
