package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func create(dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	certName, keyName := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	for _, name := range []string{certName, keyName} {
		if _, err := os.Lstat(name); err == nil {
			return fmt.Errorf("refusing to overwrite %s", name)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}
	now := time.Now()
	template := &x509.Certificate{SerialNumber: serial, Subject: pkix.Name{CommonName: "localhost"}, DNSNames: []string{"localhost"}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(24 * time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return err
	}
	encodedKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}
	// O_EXCL protects against accidental overwrite even if a file appears after Lstat.
	write := func(name string, block *pem.Block) error {
		f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return err
		}
		encodeErr := pem.Encode(f, block)
		closeErr := f.Close()
		if encodeErr != nil {
			return encodeErr
		}
		return closeErr
	}
	if err = write(keyName, &pem.Block{Type: "PRIVATE KEY", Bytes: encodedKey}); err != nil {
		return err
	}
	if err = write(certName, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return err
	}
	fmt.Printf("Created %s and %s (local development, expires in 24h)\n", certName, keyName)
	return nil
}
func main() {
	dir := flag.String("dir", ".local-certs", "directory for new certificate and key")
	flag.Parse()
	if err := create(*dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
