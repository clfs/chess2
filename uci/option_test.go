package uci

import (
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
			Option{Name: "Move Overhead", Type: SpinOptionType, Default: "10", Min: 0, Max: 5000},
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
