package lab

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type Profile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func DecodeProfile(data []byte) (Profile, error) {
	var raw struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&raw); err != nil {
		return Profile{}, err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return Profile{}, errors.New("expected exactly one JSON value")
	}
	if raw.Name == nil || strings.TrimSpace(*raw.Name) == "" || raw.Age == nil || *raw.Age < 0 {
		return Profile{}, errors.New("invalid profile")
	}
	return Profile{*raw.Name, *raw.Age}, nil
}
