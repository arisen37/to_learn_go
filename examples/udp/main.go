package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"
)

func run() error {
	server, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer server.Close()
	if err = server.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		n, addr, err := server.ReadFrom(buf)
		if err == nil {
			_, err = server.WriteTo(buf[:n], addr)
		}
		done <- err
	}()
	defer func() { server.Close(); <-done }()
	client, err := net.Dial("udp", server.LocalAddr().String())
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		return err
	}
	payload := []byte("one datagram")
	if _, err = client.Write(payload); err != nil {
		return err
	}
	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		return err
	}
	if !bytes.Equal(buf[:n], payload) {
		return fmt.Errorf("unexpected datagram")
	}
	fmt.Println(string(buf[:n]))
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
