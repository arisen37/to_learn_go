package repair

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func newTestStore() *Store { return NewStore(func() time.Time { return time.Unix(100, 0) }) }
func TestTitlesAndSnapshots(t *testing.T) {
	s := newTestStore()
	for _, title := range []string{"", " \t ", strings.Repeat("界", 121)} {
		if _, err := s.Create(title); !errors.Is(err, ErrTitle) {
			t.Fatalf("title %q: %v", title, err)
		}
	}
	a, err := s.Create("  Go界  ")
	if err != nil || a.Title != "Go界" || a.ID != 1 || !a.CreatedAt.Equal(time.Unix(100, 0)) {
		t.Fatal(a, err)
	}
	a.Title = "changed"
	saved, err := s.Get(1)
	if err != nil || saved.Title != "Go界" {
		t.Fatal("leaked pointer", saved, err)
	}
	if s.Count() != 1 {
		t.Fatal(s.Count())
	}
}
func TestMissingResources(t *testing.T) {
	s := newTestStore()
	if _, err := s.Get(9); !errors.Is(err, ErrNotFound) {
		t.Fatal("missing Get", err)
	}
	if _, err := s.Complete(9, true); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
	if err := s.Delete(9); !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}
}
func TestPaginationAndOrder(t *testing.T) {
	s := newTestStore()
	for i := 0; i < 20; i++ {
		if _, err := s.Create("task"); err != nil {
			t.Fatal(err)
		}
	}
	for attempt := 0; attempt < 20; attempt++ {
		tasks, err := s.List(2, 3)
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 3 || tasks[0].ID != 3 || tasks[2].ID != 5 {
			t.Fatal("unstable page", tasks)
		}
	}
	for _, offset := range []int{20, 1000} {
		tasks, err := s.List(offset, 10)
		if err != nil || len(tasks) != 0 || tasks == nil {
			t.Fatal(tasks, err)
		}
	}
	for _, page := range [][2]int{{-1, 10}, {0, 0}, {0, 101}} {
		if _, err := s.List(page[0], page[1]); !errors.Is(err, ErrPage) {
			t.Fatal(page, err)
		}
	}
	tasks, _ := s.List(0, 1)
	tasks[0].Title = "outside"
	saved, _ := s.Get(1)
	if saved.Title == "outside" {
		t.Fatal("snapshot mutation")
	}
}
func TestReopenAndStats(t *testing.T) {
	s := newTestStore()
	task, _ := s.Create("Go")
	s.Create("CSV")
	if _, err := s.Complete(task.ID, true); err != nil {
		t.Fatal(err)
	}
	if got := s.Stats(); got != (Summary{Total: 2, Completed: 1, Pending: 1}) {
		t.Fatal(got)
	}
	reopened, err := s.Complete(task.ID, false)
	if err != nil || reopened.Done {
		t.Fatal("cannot reopen", reopened, err)
	}
	if got := s.Stats(); got.Completed != 0 || got.Pending != 2 {
		t.Fatal(got)
	}
	if err = s.Delete(task.ID); err != nil || s.Count() != 1 {
		t.Fatal(err)
	}
}
func TestImportAtomicity(t *testing.T) {
	s := newTestStore()
	s.Create("existing")
	for _, bad := range []string{"title,done\nvalid,true\n ,false\n", "title,done\nvalid,true\nbad,maybe\n", "title,done\nvalid,true\nbad\n", "wrong,header\n"} {
		before := s.Snapshot()
		if n, err := s.ImportCSV(strings.NewReader(bad)); err == nil || n != 0 {
			t.Fatal("invalid import accepted", n, err)
		}
		if !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal("partial import changed store")
		}
	}
	n, err := s.ImportCSV(strings.NewReader("title,done\n\"learn, go\",true\nnext,false\n"))
	if err != nil || n != 2 {
		t.Fatal(n, err)
	}
	tasks := s.Snapshot()
	if tasks[1].ID != 2 || !tasks[1].Done || tasks[1].Title != "learn, go" {
		t.Fatal(tasks)
	}
}

var errSink = errors.New("sink failed")

type failedSink struct{}

func (failedSink) Write([]byte) (int, error) { return 0, errSink }
func TestExportFlushFailure(t *testing.T) {
	s := newTestStore()
	s.Create("a,b")
	if err := s.ExportCSV(failedSink{}); !errors.Is(err, errSink) {
		t.Fatal("lost flush error", err)
	}
	var out strings.Builder
	if err := s.ExportCSV(&out); err != nil {
		t.Fatal(err)
	}
	rows, err := csv.NewReader(strings.NewReader(out.String())).ReadAll()
	if err != nil || len(rows) != 2 || rows[1][1] != "a,b" {
		t.Fatal(rows, err)
	}
}
func TestSearch(t *testing.T) {
	s := newTestStore()
	s.Create("Learn Go")
	s.Create("JSON")
	s.Create("go routines")
	tasks, err := s.Search(" GO ")
	if err != nil || len(tasks) != 2 || tasks[0].ID != 1 || tasks[1].ID != 3 {
		t.Fatal(tasks, err)
	}
	if _, err = s.Search(" "); err == nil {
		t.Fatal("blank search")
	}
}
func TestHTTPBoundaries(t *testing.T) {
	s := newTestStore()
	s.Create("Go")
	h := NewHandler(s)
	for _, tc := range []struct {
		method, url, body string
		status            int
	}{
		{"GET", "/health", "", 200}, {"GET", "/tasks/99", "", 404}, {"GET", "/tasks?offset=999", "", 200}, {"GET", "/tasks?offset=oops", "", 400}, {"GET", "/tasks?limit=0", "", 400}, {"GET", "/tasks?limit=1&limit=2", "", 400},
		{"POST", "/tasks", `{"title":"next"}`, 201}, {"POST", "/tasks", `{"title":"   "}`, 400}, {"POST", "/tasks", `{"title":"x","typo":true}`, 400}, {"POST", "/tasks", `{"title":"x"} {}`, 400}, {"POST", "/tasks", `{"title":"x"} garbage`, 400}, {"POST", "/tasks", `null`, 400}, {"POST", "/tasks", strings.Repeat(" ", 1025), 413},
		{"PATCH", "/tasks/1", `{"done":true}`, 200}, {"PATCH", "/tasks/1", `{"done":false}`, 200}, {"PATCH", "/tasks/1", `{}`, 400}, {"PATCH", "/tasks/1", `{"done":true} garbage`, 400}, {"GET", "/tasks/no", "", 400}, {"PUT", "/tasks/1", "", 405}, {"GET", "/unknown", "", 404}, {"GET", "/stats", "", 200}, {"GET", "/search?q=GO", "", 200}, {"GET", "/search", "", 400},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body)))
		if w.Code != tc.status {
			t.Fatalf("%s %s %s got=%d want=%d body=%s", tc.method, tc.url, tc.body, w.Code, tc.status, w.Body.String())
		}
		if w.Code == 405 && w.Header().Get("Allow") == "" {
			t.Fatal("missing Allow")
		}
		if w.Code == 201 && w.Header().Get("Location") == "" {
			t.Fatal("missing Location")
		}
		if w.Code == 200 || w.Code == 201 {
			if !json.Valid(w.Body.Bytes()) || w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("invalid JSON response", w.Body.String())
			}
		}
	}
	task, _ := s.Get(1)
	if task.Done {
		t.Fatal("PATCH false ignored")
	}
}
func TestConcurrentStore(t *testing.T) {
	s := newTestStore()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 40; j++ {
				task, err := s.Create("concurrent")
				if err != nil {
					t.Error(err)
					return
				}
				s.Complete(task.ID, true)
				s.Get(task.ID)
				s.List(0, 10)
				s.Stats()
				s.Search("current")
				if err = s.ExportCSV(io.Discard); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
	if s.Count() != 640 {
		t.Fatal("lost writes", s.Count())
	}
	summary := s.Stats()
	if summary.Completed != 640 || summary.Pending != 0 {
		t.Fatal(summary)
	}
}

func TestOffsetPastEnd(t *testing.T) {
	// Recover only so this deliberate starter panic does not hide other regressions.
	defer func() {
		if value := recover(); value != nil {
			t.Fatalf("out-of-range offset panicked: %v", value)
		}
	}()
	s := newTestStore()
	s.Create("one")
	tasks, err := s.List(999, 10)
	if err != nil || tasks == nil || len(tasks) != 0 {
		t.Fatal(tasks, err)
	}
}

func TestExtraJSONRejected(t *testing.T) {
	s := newTestStore()
	h := NewHandler(s)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"one"} {"title":"two"}`)))
	if w.Code != 400 || s.Count() != 0 {
		t.Fatalf("extra JSON accepted: status=%d count=%d", w.Code, s.Count())
	}
}
