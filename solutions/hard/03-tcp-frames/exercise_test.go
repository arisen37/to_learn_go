package lab

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

type shortWriter struct{ bytes.Buffer }

func (w *shortWriter) Write(p []byte) (int, error) {
	if len(p) > 2 {
		p = p[:2]
	}
	return w.Buffer.Write(p)
}

type zeroWriter struct{}

func (zeroWriter) Write([]byte) (int, error) { return 0, nil }

var errWrite = errors.New("write failure")

type badWriter struct{}

func (badWriter) Write([]byte) (int, error) { return 0, errWrite }
func TestFrames(t *testing.T) {
	w := &shortWriter{}
	for _, p := range [][]byte{[]byte("hello"), nil, []byte("世界")} {
		if err := WriteFrame(w, p); err != nil {
			t.Fatal(err)
		}
	}
	for _, want := range [][]byte{[]byte("hello"), nil, []byte("世界")} {
		got, err := ReadFrame(&w.Buffer)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatal(got, err)
		}
	}
	for _, bad := range [][]byte{{0, 0}, {0, 0, 0, 3, 'a'}, {0, 16, 0, 1}} {
		if _, err := ReadFrame(bytes.NewReader(bad)); err == nil {
			t.Fatal("accepted truncated/oversized")
		}
	}
	if !errors.Is(WriteFrame(zeroWriter{}, []byte("x")), io.ErrShortWrite) {
		t.Fatal("no progress")
	}
	if !errors.Is(WriteFrame(badWriter{}, nil), errWrite) {
		t.Fatal("write error")
	}
	if err := WriteFrame(io.Discard, make([]byte, (1<<20)+1)); err == nil {
		t.Fatal("oversized write")
	}
}
func TestFramesPipe(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	_ = a.SetDeadline(time.Now().Add(time.Second))
	_ = b.SetDeadline(time.Now().Add(time.Second))
	done := make(chan error, 1)
	go func() { done <- WriteFrame(a, []byte("network")) }()
	p, err := ReadFrame(b)
	if err != nil || string(p) != "network" {
		t.Fatal(string(p), err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
