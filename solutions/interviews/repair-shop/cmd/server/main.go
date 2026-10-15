package main

import (
	"context"
	"errors"
	repair "example/hello/solutions/interviews/repair-shop"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func run(ctx context.Context, addr string) error {
	store := repair.NewStore(nil)
	if _, err := store.Create("Read the architecture before editing"); err != nil {
		return err
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           repair.NewHandler(store),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(shutdown)
		if err != nil {
			_ = server.Close()
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
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, *addr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
