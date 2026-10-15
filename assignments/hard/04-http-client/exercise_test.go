package lab

import (
	"context"
	"errors"
	"example/hello/internal/testutil"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }
func TestGetJSON(t *testing.T) {
	for _, tc := range []struct {
		status int
		body   string
		ok     bool
	}{{200, `{"n":3}`, true}, {500, `{}`, false}, {200, `{`, false}, {200, `{} {}`, false}, {200, strings.Repeat(" ", (1<<20)+1), false}} {
		listener := testutil.Listen()
		srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Error("method")
			}
			w.WriteHeader(tc.status)
			_, _ = io.WriteString(w, tc.body)
		})}
		serveDone := make(chan error, 1)
		go func() { serveDone <- srv.Serve(listener) }()
		transport := &http.Transport{DialContext: listener.DialContext}
		client := &http.Client{Transport: transport, Timeout: time.Second}
		var out struct{ N int }
		err := GetJSON(context.Background(), client, "http://local.test", &out)
		transport.CloseIdleConnections()
		_ = srv.Close()
		<-serveDone
		if (err == nil) != tc.ok {
			t.Fatalf("status=%d err=%v", tc.status, err)
		}
		if tc.ok && out.N != 3 {
			t.Fatal(out)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out any
	if err := GetJSON(ctx, http.DefaultClient, "http://127.0.0.1:1", &out); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if err := GetJSON(context.Background(), nil, "http://example.invalid", &out); err == nil {
		t.Fatal("nil client")
	}
	if err := GetJSON(context.Background(), http.DefaultClient, "://bad", &out); err == nil {
		t.Fatal("invalid URL")
	}
	for _, status := range []int{200, 500} {
		body := &trackedBody{Reader: strings.NewReader(`{}`)}
		c := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Body: body, Header: make(http.Header)}, nil
		})}
		_ = GetJSON(context.Background(), c, "http://local.test", &out)
		if !body.closed {
			t.Fatal("body not closed")
		}
	}
}
