package lab

import "testing"

func TestProfile(t *testing.T) {
	p, err := DecodeProfile([]byte(`{"name":" Go ","age":0}`))
	if err != nil || p.Name != " Go " || p.Age != 0 {
		t.Fatal(p, err)
	}
	for _, bad := range []string{`null`, `{}`, `{"name":"x"}`, `{"name":" ","age":1}`, `{"name":"x","age":-1}`, `{"name":"x","age":"1"}`, `{"name":"x","age":null}`, `{"name":"x","age":1,"extra":true}`, `{"name":"x","age":1} {}`, `{"name":"x","age":1} garbage`, `[`} {
		if _, err := DecodeProfile([]byte(bad)); err == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
}
