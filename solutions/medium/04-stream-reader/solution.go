package lab

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func SumLines(r io.Reader) (int64, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 4096), 1<<20)
	var total int64
	line := 0
	for s.Scan() {
		line++
		v := strings.TrimSpace(s.Text())
		if v == "" {
			continue
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("line %d: %w", line, err)
		}
		if n > 0 && total > math.MaxInt64-n || n < 0 && total < math.MinInt64-n {
			return 0, errors.New("sum overflow")
		}
		total += n
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return total, nil
}
