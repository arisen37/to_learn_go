package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func run(addr, message string) error {
	if strings.ContainsAny(message, "\r\n") || len(message) > 4095 {
		return fmt.Errorf("one line of at most 4095 bytes required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err = conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return err
	}
	if _, err = io.WriteString(conn, message+"\n"); err != nil {
		return err
	}
	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 1024), 4097)
	if !scanner.Scan() {
		if err = scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	fmt.Println(scanner.Text())
	return nil
}
func main() {
	addr := flag.String("addr", "127.0.0.1:3000", "server address")
	message := flag.String("message", "hello Go", "one-line message")
	flag.Parse()
	if err := run(*addr, *message); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
