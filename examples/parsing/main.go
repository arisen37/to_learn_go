package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func run() error {
	key, value, ok := strings.Cut("age=21", "=")
	if !ok {
		return fmt.Errorf("missing equals")
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	fmt.Println(key, n, strings.Fields(" Go\tJSON CSV "))
	var p Profile
	d := json.NewDecoder(strings.NewReader(`{"name":"Gopher","age":21}`))
	d.DisallowUnknownFields()
	if err = d.Decode(&p); err != nil {
		return err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	encoded, err := json.Marshal(p)
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	reader := csv.NewReader(strings.NewReader("name,score\n\"Doe, Jane\",90\n"))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	writer := csv.NewWriter(os.Stdout)
	if err = writer.WriteAll(records); err != nil {
		return err
	}
	return writer.Error()
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
