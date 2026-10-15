package lab

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPHandler(t *testing.T) {
	h := NewHandler()
	for _, tc := range []struct {
		method, url, body string
		status            int
	}{
		{"POST", "/sum", `{"a":2,"b":3}`, 200}, {"GET", "/health", "", 200}, {"GET", "/sum", "", 405}, {"POST", "/health", "", 405}, {"GET", "/missing", "", 404},
		{"POST", "/sum", `{"a":1}`, 400}, {"POST", "/sum", `null`, 400}, {"POST", "/sum", `{"a":1,"b":null}`, 400}, {"POST", "/sum", `{"a":1,"b":2,"x":0}`, 400}, {"POST", "/sum", `{"a":1,"b":2} {}`, 400}, {"POST", "/sum", `{"a":9223372036854775807,"b":1}`, 400}, {"POST", "/sum", strings.Repeat(" ", 1025), 413},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Fatalf("%s %s %q got %d want %d", tc.method, tc.url, tc.body, w.Code, tc.status)
		}
		if w.Code == 405 && w.Header().Get("Allow") == "" {
			t.Fatal("missing Allow")
		}
		if tc.url == "/sum" && w.Code == 200 {
			var obj map[string]int
			if err := json.Unmarshal(w.Body.Bytes(), &obj); err != nil || obj["sum"] != 5 || w.Header().Get("Content-Type") != "application/json" {
				t.Fatal(w.Body.String(), err)
			}
		}
	}
}
