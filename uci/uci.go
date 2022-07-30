// Package uci implements a client for the Universal Chess Interface (UCI)
// protocol.
package uci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
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
	Min     string
	Max     string
	Vars    []string
}

// UnmarshalText unmarshals a textual representation of an Option.
func (o *Option) UnmarshalText(text []byte) error {
	fields := bytes.Fields(text)

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

	var nameBuf bytes.Buffer
	for ; pos < len(fields); pos++ {
		if string(fields[pos]) == "type" {
			break
		}
		nameBuf.Write(fields[pos])
		nameBuf.WriteByte(' ')
	}
	o.Name = strings.TrimSpace(nameBuf.String())

	for ; pos < len(fields)-1; pos++ {
		cur, nxt := fields[pos], fields[pos+1]
		switch string(cur) {
		case "type":
			o.Type = string(nxt)
		case "default":
			o.Default = string(nxt)
		case "min":
			o.Min = string(nxt)
		case "max":
			o.Max = string(nxt)
		case "var":
			o.Vars = append(o.Vars, string(nxt))
		}
	}
	return nil
}

// Client is a UCI-compatible client.
type Client struct {
	r io.Reader
	w io.Writer

	Name    string   // The name of the engine.
	Author  string   // The author of the engine.
	Options []Option // Available engine options.
	Result  Result   // The result of the last search.
}

type Result struct {
	// todo
}

// NewClient returns a client that reads UCI responses from r and writes UCI
// commands to w.
func NewClient(r io.Reader, w io.Writer) *Client {
	return &Client{r: r, w: w}
}

// NewClientFromPath runs the engine located at path and returns a client
// connected to the engine.
func NewClientFromPath(path string) (*Client, error) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return NewClient(stdout, stdin), nil
}

func (c *Client) send(s string) error {
	_, err := c.w.Write([]byte(s + "\n"))
	return err
}

// UCI sends the "uci" command. It tells the engine to use the UCI protocol.
func (c *Client) UCI() error {
	if err := c.send("uci"); err != nil {
		return err
	}

	s := bufio.NewScanner(c.r)
outer:
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "id name "):
			c.Name = strings.TrimPrefix(line, "id name ")
		case strings.HasPrefix(line, "id author "):
			c.Author = strings.TrimPrefix(line, "id author ")
		case strings.HasPrefix(line, "option "):
			var opt Option
			if err := opt.UnmarshalText([]byte(line)); err != nil {
				return err
			}
			c.Options = append(c.Options, opt)
		case line == "uciok":
			break outer
		}
	}

	return s.Err()
}

// Debug sends the "debug" command. It toggles the engine's debug mode. Many
// engines do not support this command.
func (c *Client) Debug(on bool) error {
	if on {
		return c.send("debug on")
	}
	return c.send("debug off")
}

// IsReady sends the "isready" command. It returns when the engine is ready to
// accept commands.
func (c *Client) IsReady() error {
	if err := c.send("isready"); err != nil {
		return err
	}

	s := bufio.NewScanner(c.r)
	for s.Scan() {
		if s.Text() == "readyok" {
			return nil
		}
	}

	return s.Err()
}

// SetOption sends the "setoption" command. It sets an option in the engine's
// internal parameters.
func (c *Client) SetOption(name, value string) error {
	return nil // todo
}

// RegisterParams contains parameters for the "register" command.
type RegisterParams struct {
	Later bool
	Name  string
	Code  string
}

// Register sends the "register" command. It submits registration information
// for licensing. Many engines do not support this command.
func (c *Client) Register(r RegisterParams) error {
	if r.Later {
		return c.send("register later")
	}
	return c.send(fmt.Sprintf("register name %s code %s", r.Name, r.Code))
}

// UCINewGame sends the "ucinewgame" command.
func (c *Client) UCINewGame() error {
	return c.send("ucinewgame")
}

// PositionParams contains parameters for the "position" command.
type PositionParams struct {
	// todo
}

// Position sends the "position" command. It sets the board position.
func (c *Client) Position(p PositionParams) error {
	return nil // todo
}

// GoParams contains parameters for the "go" command.
type GoParams struct {
	// todo
}

// Go sends the "go" command. It starts engine calculations.
func (c *Client) Go(p GoParams) error {
	return nil // todo
}

// Stop sends the "stop" command. It stops engine calculations.
func (c *Client) Stop() error {
	return c.send("stop")
}

// PonderHit sends the "ponderhit" command. It tells the engine that the
// opponent has played its best move.
func (c *Client) PonderHit() error {
	return c.send("ponderhit")
}

// Quit sends the "quit" command. It tells the engine to quit.
func (c *Client) Quit() error {
	return c.send("quit")
}
