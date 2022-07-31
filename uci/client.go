package uci

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

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

// NewClient returns a UCI client that reads from r and writes to w.
func NewClient(r io.Reader, w io.Writer) *Client {
	return &Client{r: r, w: w}
}

// NewClientFromPath runs the engine located at path and returns a client
// connected to the engine's standard input and output.
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

// UCI sends the "uci" command. It tells the engine to use the UCI protocol.
func (c *Client) UCI() error {
	fmt.Fprintln(c.w, "uci")

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

// Debug sends the "debug" command. It toggles the engine's debug mode.
func (c *Client) Debug(on bool) {
	if on {
		fmt.Fprintln(c.w, "debug on")
	}
	fmt.Fprintln(c.w, "debug off")
}

// IsReady sends the "isready" command. It blocks until the engine is ready to
// accept commands.
func (c *Client) IsReady() error {
	fmt.Fprintln(c.w, "isready")

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
// for licensing.
func (c *Client) Register(r RegisterParams) {
	if r.Later {
		fmt.Fprintln(c.w, "register later")
	}
	fmt.Fprintf(c.w, "register name %s code %s\n", r.Name, r.Code)
}

// UCINewGame sends the "ucinewgame" command. This tells the engine the next
// search will be from a different game.
func (c *Client) UCINewGame() {
	fmt.Fprintln(c.w, "ucinewgame")
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
func (c *Client) Stop() {
	fmt.Fprintln(c.w, "stop")
}

// PonderHit sends the "ponderhit" command. It tells the engine that the
// opponent has played its best move.
func (c *Client) PonderHit() {
	fmt.Fprintln(c.w, "ponderhit")
}

// Quit sends the "quit" command. It tells the engine to quit.
func (c *Client) Quit() {
	fmt.Fprintln(c.w, "quit")
}
