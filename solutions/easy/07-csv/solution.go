package lab

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Score struct {
	Name  string
	Value int
}

func ReadScores(r io.Reader) ([]Score, error) {
	c := csv.NewReader(r)
	header, err := c.Read()
	if err != nil {
		return nil, err
	}
	if len(header) != 2 || header[0] != "name" || header[1] != "score" {
		return nil, errors.New("invalid header")
	}
	out := make([]Score, 0)
	for {
		row, err := c.Read()
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		n, err := strconv.Atoi(row[1])
		if err != nil || n < 0 || n > 100 || strings.TrimSpace(row[0]) == "" {
			return nil, errors.New("invalid score record")
		}
		out = append(out, Score{row[0], n})
	}
}
