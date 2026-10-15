package repair

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// NewHandler adapts HTTP to the store, keeping wire validation at the edge.
func NewHandler(store *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/health":
			if !allow(w, r, "GET") {
				return
			}
			writeJSON(w, 200, map[string]bool{"ok": true})
		case r.URL.Path == "/stats":
			if !allow(w, r, "GET") {
				return
			}
			writeJSON(w, 200, store.Stats())
		case r.URL.Path == "/search":
			if !allow(w, r, "GET") {
				return
			}
			values := r.URL.Query()["q"]
			if len(values) != 1 {
				http.Error(w, "one search query required", 400)
				return
			}
			tasks, err := store.Search(values[0])
			if err != nil {
				http.Error(w, "invalid search", 400)
				return
			}
			writeJSON(w, 200, tasks)
		case r.URL.Path == "/tasks":
			handleCollection(w, r, store)
		case r.URL.Path == "/export":
			if !allow(w, r, "GET") {
				return
			}
			// Buffer before committing headers, so export failure can return 500.
			var buf bytes.Buffer
			if err := store.ExportCSV(&buf); err != nil {
				http.Error(w, "export failed", 500)
				return
			}
			w.Header().Set("Content-Type", "text/csv")
			_, _ = w.Write(buf.Bytes())
		case strings.HasPrefix(r.URL.Path, "/tasks/"):
			raw := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(raw)
			if err != nil || id <= 0 {
				http.Error(w, "invalid id", 400)
				return
			}
			handleTask(w, r, store, id)
		default:
			http.NotFound(w, r)
		}
	})
}

func handleCollection(w http.ResponseWriter, r *http.Request, s *Store) {
	switch r.Method {
	case "GET":
		offset, err := queryInt(r, "offset", 0)
		if err != nil {
			http.Error(w, "invalid offset", 400)
			return
		}
		limit, err := queryInt(r, "limit", 10)
		if err != nil {
			http.Error(w, "invalid limit", 400)
			return
		}
		tasks, err := s.List(offset, limit)
		if err != nil {
			http.Error(w, "invalid page", 400)
			return
		}
		writeJSON(w, 200, tasks)
	case "POST":
		var body struct {
			Title *string `json:"title"`
		}
		if !decode(w, r, &body) {
			return
		}
		if body.Title == nil {
			http.Error(w, "title required", 400)
			return
		}
		task, err := s.Create(*body.Title)
		if err != nil {
			http.Error(w, "invalid title", 400)
			return
		}
		w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
		writeJSON(w, 201, task)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", 405)
	}
}

func handleTask(w http.ResponseWriter, r *http.Request, s *Store, id int) {
	switch r.Method {
	case "GET":
		task, err := s.Get(id)
		if err != nil {
			storeError(w, err)
			return
		}
		writeJSON(w, 200, task)
	case "PATCH":
		var body struct {
			Done *bool `json:"done"`
		}
		if !decode(w, r, &body) {
			return
		}
		if body.Done == nil {
			http.Error(w, "done required", 400)
			return
		}
		task, err := s.Complete(id, *body.Done)
		if err != nil {
			storeError(w, err)
			return
		}
		writeJSON(w, 200, task)
	case "DELETE":
		if err := s.Delete(id); err != nil {
			storeError(w, err)
			return
		}
		w.WriteHeader(204)
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		http.Error(w, "method not allowed", 405)
	}
}

func queryInt(r *http.Request, key string, fallback int) (int, error) {
	values, ok := r.URL.Query()[key]
	if !ok {
		return fallback, nil
	}
	if len(values) != 1 || values[0] == "" {
		return 0, ErrPage
	}
	return strconv.Atoi(values[0])
}

func decode(w http.ResponseWriter, r *http.Request, out any) bool {
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024))
	if err != nil {
		var limit *http.MaxBytesError
		if errors.As(err, &limit) {
			http.Error(w, "body too large", 413)
		} else {
			http.Error(w, "body read failed", 400)
		}
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(out); err != nil {
		http.Error(w, "invalid JSON", 400)
		return false
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		http.Error(w, "trailing JSON", 400)
		return false
	}
	return true
}

func allow(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	http.Error(w, "method not allowed", 405)
	return false
}

func storeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", 404)
		return
	}
	http.Error(w, "internal error", 500)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "encode failed", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(data, '\n'))
}
