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
		name    string
		author  string
		options []Option
	}{
		"My Chess Engine",
		"Firstname Lastname",
		[]Option{
			{Name: "DoFoo", Type: ButtonOptionType, Default: ""},
			{Name: "Fruit", Type: ComboOptionType, Default: "apple", Vars: []string{"apple", "banana"}},
		},
	}
	c := NewClient(bytes.NewReader(data), io.Discard)
	if err := c.UCI(); err != nil {
		t.Errorf("err: %v", err)
	}
	if want.name != c.Name {
		t.Errorf("name: want %s, got %s", want.name, c.Name)
	}
	if want.author != c.Author {
		t.Errorf("author: want %s, got %s", want.author, c.Author)
	}
	if diff := cmp.Diff(want.options, c.Options, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
		t.Errorf("options: mismatch (-want +got):\n%s", diff)
	}
}
