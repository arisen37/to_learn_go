package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	mux.HandleFunc("GET /time", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"time": time.Now().UTC().Format(time.RFC3339)})
	})
	return mux
}
func run(ctx context.Context, addr, cert, key string) error {
	srv := &http.Server{Addr: addr, Handler: handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	done := make(chan error, 1)
	go func() { done <- srv.ListenAndServeTLS(cert, key) }()
	slog.Info("HTTPS server starting", "address", addr)
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		stop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := srv.Shutdown(stop)
		if err != nil {
			_ = srv.Close()
		}
		serveErr := <-done
		if err != nil {
			return err
		}
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			return serveErr
		}
		return nil
	}
}
func main() {
	addr := flag.String("addr", "127.0.0.1:8443", "listen address")
	cert := flag.String("cert", ".local-certs/cert.pem", "PEM certificate")
	key := flag.String("key", ".local-certs/key.pem", "PEM private key")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *addr, *cert, *key); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
