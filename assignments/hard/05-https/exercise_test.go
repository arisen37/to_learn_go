package lab

import (
	"context"
	"crypto/tls"
	"example/hello/internal/testutil"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServerConfig(t *testing.T) {
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	s := NewServer("127.0.0.1:0", h)
	if s.Addr != "127.0.0.1:0" || s.Handler == nil || s.ReadHeaderTimeout != 5*time.Second || s.ReadTimeout != 10*time.Second || s.WriteTimeout != 15*time.Second || s.IdleTimeout != 60*time.Second || s.MaxHeaderBytes != 1<<20 || s.TLSConfig == nil || s.TLSConfig.MinVersion < tls.VersionTLS12 {
		t.Fatalf("bad server config: %+v", s)
	}
}
func TestServeStartupError(t *testing.T) {
	l := testutil.Listen()
	if err := Serve(context.Background(), NewServer("", http.NotFoundHandler()), l, "missing-cert", "missing-key"); err == nil {
		t.Fatal("certificate error hidden")
	}
	if _, err := l.Accept(); err != net.ErrClosed {
		t.Fatal("listener left open", err)
	}
}
func TestServeTLSAndCancel(t *testing.T) {
	certFile, keyFile, roots, err := testutil.Certificate(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	l := testutil.Listen()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewServer("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "secure") }))
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, s, l, certFile, keyFile) }()
	transport := &http.Transport{DialContext: l.DialContext, TLSClientConfig: testutil.TLSConfig(roots)}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: time.Second}
	resp, err := client.Get("https://127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil || string(body) != "secure" || resp.TLS == nil {
		t.Fatal(string(body), err)
	}
	transport.CloseIdleConnections()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop")
	}
}
