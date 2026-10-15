package lab

import (
	"context"
	"io"
)

func Aggregate(ctx context.Context, r io.Reader, workers int) (map[string]int64, error) {
	panic("TODO")
}
