package lab

import (
	"context"
	"errors"
	"io"
	"maps"
	"strings"
	"testing"
)

var errSource = errors.New("source failed")

type sourceFailure struct{}

func (sourceFailure) Read([]byte) (int, error) { return 0, errSource }
func TestAggregate(t *testing.T) {
	for _, workers := range []int{1, 4} {
		got, err := Aggregate(context.Background(), strings.NewReader("key,value\n a ,2\nb,5\na,-1\n\"x,y\",3\n"), workers)
		if err != nil || !maps.Equal(got, map[string]int64{"a": 1, "b": 5, "x,y": 3}) {
			t.Fatal(got, err)
		}
	}
	for _, bad := range []string{"bad,header\n", "key,value\na\n", "key,value\na,no\n", "key,value\n ,3\n", "key,value\na,9223372036854775807\na,1\n"} {
		got, err := Aggregate(context.Background(), strings.NewReader(bad), 2)
		if err == nil || got != nil {
			t.Fatalf("bad=%q got=%v err=%v", bad, got, err)
		}
	}
	for _, empty := range []string{"", "key,value\n"} {
		got, err := Aggregate(context.Background(), strings.NewReader(empty), 2)
		if err != nil || len(got) != 0 {
			t.Fatal(got, err)
		}
	}
	if _, err := Aggregate(context.Background(), strings.NewReader(""), 0); err == nil {
		t.Fatal("workers")
	}
	if _, err := Aggregate(context.Background(), io.MultiReader(strings.NewReader("key,value\na,1\n"), sourceFailure{}), 2); !errors.Is(err, errSource) {
		t.Fatal("read error", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Aggregate(ctx, strings.NewReader("key,value\na,1\n"), 2); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
