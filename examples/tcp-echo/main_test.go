package main

import (
	"bufio"
	"context"
	"example/hello/internal/testutil"
	"io"
	"testing"
	"time"
)

func TestServeAndCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := testutil.Listen()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, l) }()
	c, err := l.DialContext(ctx, "tcp", "")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(time.Second))
	if _, err = io.WriteString(c, "hello\n"); err != nil {
		t.Fatal(err)
	}
	line, err := bufio.NewReader(c).ReadString('\n')
	if err != nil || line != "hello\n" {
		t.Fatal(line, err)
	}
	cancel()
	select {
	case err = <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not join handlers")
	}
}
