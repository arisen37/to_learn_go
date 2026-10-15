package lab

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadScores(t *testing.T) {
	got, err := ReadScores(strings.NewReader("name,score\n\"Doe, Jane\",100\nGo,0\n"))
	want := []Score{{"Doe, Jane", 100}, {"Go", 0}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatal(got, err)
	}
	got, err = ReadScores(strings.NewReader("name,score\n\"Multi\nLine\",8\n"))
	if err != nil || len(got) != 1 || got[0].Name != "Multi\nLine" {
		t.Fatal(got, err)
	}
	for _, bad := range []string{"", "score,name\n", "name,score\nx,101\n", "name,score\nx,-1\n", "name,score\nx,no\n", "name,score\nx\n", "name,score\n ,5\n", "name,score\nx,5\ny,no\n"} {
		got, err := ReadScores(strings.NewReader(bad))
		if err == nil || got != nil {
			t.Fatalf("accepted/partial result for %q", bad)
		}
	}
	if got, err := ReadScores(strings.NewReader("name,score\n")); err != nil || len(got) != 0 {
		t.Fatal(got, err)
	}
}
