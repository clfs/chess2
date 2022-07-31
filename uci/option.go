package uci

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CheckOptionType  = "check"  // A boolean option.
	SpinOptionType   = "spin"   // An integer option in a certain range.
	ComboOptionType  = "combo"  // A string option from a list of available strings.
	ButtonOptionType = "button" // A button option that causes an effect when set.
	StringOptionType = "string" // A string option.
)

// Option represents an option that engines can set.
type Option struct {
	Name    string
	Type    string
	Default string
	Min     int
	Max     int
	Vars    []string
}

func (o *Option) UnmarshalText(text []byte) error {
	fields := strings.Fields(string(text))

	if len(fields) < 5 {
		return fmt.Errorf("todo")
	}

	var pos int

	if f := fields[pos]; string(f) != "option" {
		return fmt.Errorf("todo")
	}
	pos++

	if f := fields[pos]; string(f) != "name" {
		return fmt.Errorf("todo")
	}
	pos++

	var acc []string
	for ; pos < len(fields); pos++ {
		if string(fields[pos]) == "type" {
			break
		}
		acc = append(acc, (fields[pos]))
	}
	o.Name = strings.Join(acc, " ")

	for ; pos < len(fields)-1; pos++ {
		cur, nxt := fields[pos], fields[pos+1]
		switch cur {
		case "type":
			o.Type = nxt
		case "default":
			o.Default = nxt
		case "min":
			min, err := strconv.Atoi(nxt)
			if err != nil {
				return fmt.Errorf("todo")
			}
			o.Min = min
		case "max":
			max, err := strconv.Atoi(nxt)
			if err != nil {
				return fmt.Errorf("todo")
			}
			o.Max = max
		case "var":
			o.Vars = append(o.Vars, nxt)
		}
	}
	return nil
}
