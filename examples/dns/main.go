package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	host := "localhost"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	addresses, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, addr := range addresses {
		fmt.Println(addr)
	}
}
