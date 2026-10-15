package lab

import "io"

type Score struct {
	Name  string
	Value int
}

func ReadScores(r io.Reader) ([]Score, error) { panic("TODO: ReadScores") }
