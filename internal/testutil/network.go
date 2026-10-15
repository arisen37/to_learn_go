// Package testutil provides in-memory networking for acceptance tests.
// It contains no assignment implementations.
package testutil

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Listener runs a real net/http or TLS stack over net.Pipe without opening ports.
type Listener struct {
	connections chan net.Conn
	done        chan struct{}
	once        sync.Once
}

func Listen() *Listener {
	return &Listener{connections: make(chan net.Conn), done: make(chan struct{})}
}
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case c := <-l.connections:
		return c, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}
func (l *Listener) Close() error   { l.once.Do(func() { close(l.done) }); return nil }
func (l *Listener) Addr() net.Addr { return memoryAddr("127.0.0.1:443") }
func (l *Listener) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	client, server := net.Pipe()
	select {
	case l.connections <- server:
		return client, nil
	case <-ctx.Done():
		client.Close()
		server.Close()
		return nil, ctx.Err()
	case <-l.done:
		client.Close()
		server.Close()
		return nil, net.ErrClosed
	}
}

type memoryAddr string

func (memoryAddr) Network() string  { return "tcp" }
func (a memoryAddr) String() string { return string(a) }

// Certificate creates a short-lived localhost certificate for tests only.
func Certificate(dir string) (certFile, keyFile string, roots *x509.CertPool, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", "", nil, err
	}
	template := &x509.Certificate{SerialNumber: serial, Subject: pkix.Name{CommonName: "localhost"}, NotBefore: time.Now().Add(-time.Minute), NotAfter: time.Now().Add(time.Hour), DNSNames: []string{"localhost"}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return "", "", nil, err
	}
	encodedKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", "", nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: encodedKey})
	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	if err = os.WriteFile(certFile, certPEM, 0600); err != nil {
		return "", "", nil, err
	}
	if err = os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		return "", "", nil, err
	}
	roots = x509.NewCertPool()
	roots.AppendCertsFromPEM(certPEM)
	return certFile, keyFile, roots, nil
}

// TLSConfig is a client config which trusts only the generated test certificate.
func TLSConfig(roots *x509.CertPool) *tls.Config {
	return &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
}
