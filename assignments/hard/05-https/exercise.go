package lab

import (
	"context"
	"net"
	"net/http"
)

func NewServer(addr string, handler http.Handler) *http.Server { panic("TODO") }
func Serve(ctx context.Context, srv *http.Server, listener net.Listener, certFile, keyFile string) error {
	panic("TODO")
}
