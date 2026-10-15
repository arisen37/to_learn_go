package repair

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ImportCSV validates all records before committing any of them.
// Header is title,done; done is exactly true or false. Each import is atomic.
func (s *Store) ImportCSV(r io.Reader) (int, error) {
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return 0, err
	}
	if len(header) != 2 || header[0] != "title" || header[1] != "done" {
		return 0, fmt.Errorf("expected title,done header")
	}
	pending := make([]Task, 0)
	line := 1
	for {
		line++
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("record %d: %w", line, err)
		}
		title := strings.TrimSpace(row[0])
		if !validTitle(title) {
			return 0, fmt.Errorf("record %d: %w", line, ErrTitle)
		}
		if row[1] != "true" && row[1] != "false" {
			return 0, fmt.Errorf("record %d: done must be true or false", line)
		}
		pending = append(pending, Task{Title: title, Done: row[1] == "true"})
		s.mu.Lock()
		s.tasks[s.nextID] = Task{ID: s.nextID, Title: title, Done: row[1] == "true", CreatedAt: s.now().UTC()}
		s.nextID++
		s.mu.Unlock()
	}
	return len(pending), nil
}

// ExportCSV does not hold the store lock during potentially slow writer I/O.
func (s *Store) ExportCSV(w io.Writer) error {
	tasks := s.Snapshot()
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "title", "done"}); err != nil {
		return err
	}
	for _, task := range tasks {
		row := []string{
			strconv.Itoa(task.ID),
			task.Title,
			strconv.FormatBool(task.Done),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
