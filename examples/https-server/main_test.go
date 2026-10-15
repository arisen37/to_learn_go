package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	for _, url := range []string{"/health", "/time"} {
		w := httptest.NewRecorder()
		handler().ServeHTTP(w, httptest.NewRequest("GET", url, nil))
		if w.Code != 200 || !json.Valid(w.Body.Bytes()) {
			t.Fatal(w.Code, w.Body.String())
		}
	}
}
