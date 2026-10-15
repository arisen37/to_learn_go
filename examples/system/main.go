package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func run() error {
	f, err := os.CreateTemp("", "go-learn-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err = io.WriteString(f, "Go uses operating-system resources through portable APIs.\n"); err != nil {
		return err
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	b, err := io.ReadAll(io.LimitReader(f, 4096))
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	fmt.Printf("temporary file: %d bytes\n%s", info.Size(), b)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "go", "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("go version: %w (%s)", err, out)
	}
	fmt.Print(string(out))
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
