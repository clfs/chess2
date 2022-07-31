package uci

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func readTestdata(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	return data
}

func TestClient_UCI(t *testing.T) {
	data := readTestdata(t, "uci-response.txt")
	want := struct {
		name   string
		author string
		opts   []Option
	}{
		"My Chess Engine",
		"Firstname Lastname",
		[]Option{
			{Name: "DoFoo", Type: ButtonOptionType, Default: ""},
			{Name: "Fruit", Type: ComboOptionType, Default: "apple", Vars: []string{"apple", "banana"}},
		},
	}
	c := NewClient(bytes.NewReader(data), io.Discard)
	name, author, opts, err := c.UCI()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if want.name != name {
		t.Errorf("name: want %s, got %s", want.name, name)
	}
	if want.author != author {
		t.Errorf("author: want %s, got %s", want.author, author)
	}
	if diff := cmp.Diff(want.opts, opts, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
		t.Errorf("opts: mismatch (-want +got):\n%s", diff)
	}
}
