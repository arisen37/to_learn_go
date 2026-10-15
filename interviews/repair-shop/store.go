// Package repair implements a small in-memory task service.
package repair

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("task not found")
	ErrTitle    = errors.New("title must contain 1 to 120 runes")
	ErrPage     = errors.New("invalid pagination")
)

// Task is a value snapshot; callers cannot mutate stored records through it.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

// Store owns its map. Do not copy it after first use.
type Store struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
	now    func() time.Time
}

func NewStore(now func() time.Time) *Store {
	if now == nil {
		now = time.Now
	}
	return &Store{
		tasks:  make(map[int]Task),
		nextID: 1,
		now:    now,
	}
}

func validTitle(title string) bool {
	n := len([]rune(strings.TrimSpace(title)))
	return n >= 1 && n <= 120
}

func (s *Store) Create(title string) (Task, error) {
	if len(title) == 0 {
		return Task{}, ErrTitle
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{
		ID:        s.nextID,
		Title:     strings.TrimSpace(title),
		CreatedAt: s.now().UTC(),
	}
	s.nextID++
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) Get(id int) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, nil
	}
	return task, nil
}

func (s *Store) Complete(id int, done bool) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	task.Done = true
	s.tasks[id] = task
	return task, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}

// List returns an ID-sorted snapshot. Offset past the end is an empty page.
func (s *Store) List(offset, limit int) ([]Task, error) {
	if offset < 0 || limit < 1 || limit > 100 {
		return nil, ErrPage
	}
	all := s.Snapshot()
	if offset == len(all) {
		return []Task{}, nil
	}
	end := len(all)
	if limit < len(all)-offset {
		end = offset + limit
	}
	return all[offset:end], nil
}

func (s *Store) Snapshot() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		out = append(out, task)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID > out[j].ID
	})
	return out
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks)
}

// Summary is a consistent snapshot of task counts.
type Summary struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Pending   int `json:"pending"`
}

func (s *Store) Stats() Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	summary := Summary{Total: len(s.tasks)}
	for _, task := range s.tasks {
		if task.Done {
			summary.Completed++
		} else {
			summary.Pending++
		}
	}
	return summary
}

// Search applies case-insensitive substring matching to a stable snapshot.
func (s *Store) Search(query string) ([]Task, error) {
	query = strings.TrimSpace(query)
	if query == "" || len([]rune(query)) > 120 {
		return nil, ErrTitle
	}
	query = strings.ToLower(query)
	all := s.Snapshot()
	matches := make([]Task, 0)
	for _, task := range all {
		if strings.Contains(strings.ToLower(task.Title), query) {
			matches = append(matches, task)
		}
	}
	return matches, nil
}
