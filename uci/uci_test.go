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

func TestOption_UnmarshalText(t *testing.T) {
	cases := []struct {
		in   []byte
		want Option
	}{
		{
			[]byte("option name WeightsFile type string default <autodiscover>"),
			Option{Name: "WeightsFile", Type: StringOptionType, Default: "<autodiscover>"},
		},
		{
			[]byte("option name BackendOptions type string default"),
			Option{Name: "BackendOptions", Type: StringOptionType},
		},
		{
			[]byte("option name Move Overhead type spin default 10 min 0 max 5000"),
			Option{Name: "Move Overhead", Type: SpinOptionType, Default: "10", Min: "0", Max: "5000"},
		},
		{
			[]byte("option name HistoryFill type combo default fen_only var no var fen_only var always"),
			Option{Name: "HistoryFill", Type: ComboOptionType, Default: "fen_only", Vars: []string{"no", "fen_only", "always"}},
		},
		{
			[]byte("option name Clear Hash type button"),
			Option{Name: "Clear Hash", Type: ButtonOptionType},
		},
		{
			[]byte("option name Ponder type check default false"),
			Option{Name: "Ponder", Type: CheckOptionType, Default: "false"},
		},
	}
	for i, c := range cases {
		var opt Option
		if err := opt.UnmarshalText(c.in); err != nil {
			t.Errorf("#%d: Option.UnmarshalText: %v", i, err)
		}
		if diff := cmp.Diff(c.want, opt, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
			t.Errorf("#%d: mismatch (-want +got):\n%s", i, diff)
		}
	}
}

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
	if diff := cmp.Diff(want.name, c.Name); diff != "" {
		t.Errorf("name: mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(want.author, c.Author); diff != "" {
		t.Errorf("author: mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(want.options, c.Options, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
		t.Errorf("options: mismatch (-want +got):\n%s", diff)
	}
}
