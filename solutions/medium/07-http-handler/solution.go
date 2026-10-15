package lab

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
)

func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			if r.Method != "GET" {
				w.Header().Set("Allow", "GET")
				http.Error(w, "method", 405)
				return
			}
			w.WriteHeader(200)
			return
		}
		if r.URL.Path != "/sum" {
			http.NotFound(w, r)
			return
		}
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method", 405)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024))
		if err != nil {
			var limit *http.MaxBytesError
			if errors.As(err, &limit) {
				http.Error(w, "too large", 413)
			} else {
				http.Error(w, "read error", 400)
			}
			return
		}
		var input struct {
			A *int64 `json:"a"`
			B *int64 `json:"b"`
		}
		d := json.NewDecoder(bytes.NewReader(body))
		d.DisallowUnknownFields()
		if err = d.Decode(&input); err != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		var extra any
		if d.Decode(&extra) != io.EOF || input.A == nil || input.B == nil {
			http.Error(w, "one object with a and b required", 400)
			return
		}
		a, b := *input.A, *input.B
		if b > 0 && a > math.MaxInt64-b || b < 0 && a < math.MinInt64-b {
			http.Error(w, "overflow", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(struct {
			Sum int64 `json:"sum"`
		}{a + b})
	})
}
