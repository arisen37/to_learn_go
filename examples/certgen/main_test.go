package main

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCertificate(t *testing.T) {
	dir := t.TempDir()
	if err := create(dir); err != nil {
		t.Fatal(err)
	}
	pair, err := tls.LoadX509KeyPair(filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(pair.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, host := range []string{"localhost", "127.0.0.1", "::1"} {
		if err = cert.VerifyHostname(host); err != nil {
			t.Fatal(err)
		}
	}
	if time.Until(cert.NotAfter) > 25*time.Hour {
		t.Fatal("expiry")
	}
	info, err := os.Stat(filepath.Join(dir, "key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0077 != 0 {
		t.Fatal("key permissions")
	}
	if err = create(dir); err == nil {
		t.Fatal("overwrote certificate")
	}
}
