package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func echo(conn net.Conn) error {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 1024), 4097)
	for {
		if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
			return err
		}
		if !scanner.Scan() {
			return scanner.Err()
		}
		line := scanner.Bytes()
		if len(line) > 4095 {
			return fmt.Errorf("line exceeds 4095 bytes")
		}
		if _, err := io.WriteString(conn, string(line)+"\n"); err != nil {
			return err
		}
	}
}
func serve(ctx context.Context, listener net.Listener) error {
	defer listener.Close()
	slots := make(chan struct{}, 32)
	var wg sync.WaitGroup
	var mu sync.Mutex
	active := map[net.Conn]struct{}{}
	watchDone := make(chan struct{})
	watchJoined := make(chan struct{})
	go func() {
		defer close(watchJoined)
		select {
		case <-ctx.Done():
			listener.Close()
			mu.Lock()
			for c := range active {
				c.Close()
			}
			mu.Unlock()
		case <-watchDone:
		}
	}()
	defer func() {
		close(watchDone)
		listener.Close()
		mu.Lock()
		for c := range active {
			c.Close()
		}
		mu.Unlock()
		wg.Wait()
		<-watchJoined
	}()
	for {
		select {
		case slots <- struct{}{}:
		case <-ctx.Done():
			return nil
		}
		conn, err := listener.Accept()
		if err != nil {
			<-slots
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		mu.Lock()
		if ctx.Err() != nil {
			mu.Unlock()
			conn.Close()
			<-slots
			return nil
		}
		active[conn] = struct{}{}
		mu.Unlock()
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			defer func() { mu.Lock(); delete(active, c); mu.Unlock(); <-slots }()
			if err := echo(c); err != nil && ctx.Err() == nil {
				fmt.Fprintln(os.Stderr, "connection:", err)
			}
		}(conn)
	}
}
func main() {
	addr := flag.String("addr", "127.0.0.1:3000", "listen address")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	listener, err := net.Listen("tcp", *addr)
	if err == nil {
		fmt.Println("listening:", listener.Addr())
		err = serve(ctx, listener)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
