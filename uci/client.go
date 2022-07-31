package uci

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
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

// UCI sends a "uci" command. It tells the engine to use the UCI protocol and
// blocks until the engine confirms.
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

// Debug sends a "debug" command. It toggles the engine's debug mode.
func (c *Client) Debug(on bool) {
	if on {
		fmt.Fprintln(c.w, "debug on")
	}
	fmt.Fprintln(c.w, "debug off")
}

// IsReady sends an "isready" command. It blocks until the engine is ready to
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

// SetOption sends a "setoption" command. It sets an option in the engine's
// internal parameters. To set a value-less option, use the empty string.
func (c *Client) SetOption(name, value string) {
	if value == "" {
		fmt.Fprintf(c.w, "setoption name %s", name)
	} else {
		fmt.Fprintf(c.w, "setoption name %s value %s", name, value)
	}
}

// Register sends a "register" command. It registers client information with the
// engine.
func (c *Client) Register(name, code string) {
	fmt.Fprintf(c.w, "register name %s code %s\n", name, code)
}

// RegisterLater sends a "register later" command. It claims that the client
// will register itself later.
func (c *Client) RegisterLater() {
	fmt.Fprintln(c.w, "register later")
}

// UCINewGame sends a "ucinewgame" command. It indicates that the next search
// will be from a different game.
func (c *Client) UCINewGame() {
	fmt.Fprintln(c.w, "ucinewgame")
}

// PositionFEN sends a "position fen" command. It sets the current position
// based on a FEN string and subsequent moves.
func (c *Client) PositionFEN(fen string, moves []string) {
	fmt.Fprintf(c.w, "position fen %s", fen)
	if len(moves) > 0 {
		fmt.Fprintf(c.w, " moves %s", strings.Join(moves, " "))
	}
	fmt.Fprintf(c.w, "\n")
}

// PositionStartPos sends a "position startpos" command. It sets the current
// position based on the standard starting position and subsequent moves.
func (c *Client) PositionStartPos(moves []string) {
	fmt.Fprintln(c.w, "position startpos")
	if len(moves) > 0 {
		fmt.Fprintf(c.w, " moves %s", strings.Join(moves, " "))
	}
	fmt.Fprintf(c.w, "\n")
}

// GoParameters contains parameters for the "go" command. Note that fields of
// type time.Duration are truncated to the millisecond.
type GoParameters struct {
	SearchMoves []string // Restrict search to these moves, if any.

	Ponder   bool          // Search in ponder mode.
	Infinite bool          // Search indefinitely.
	Mate     int           // Search for a mate in this many moves. 0 is ignored.
	MoveTime time.Duration // Search for this long. 0 is ignored.

	WhiteTime      time.Duration // Time remaining for White. 0 is infinite.
	BlackTime      time.Duration // Time remaining for Black. 0 is infinite.
	WhiteIncrement time.Duration // Time increment for White. 0 is no increment.
	BlackIncrement time.Duration // Time increment for Black. 0 is no increment.
	MovesToGo      int           // Moves remaining until next time control. 0 is ignored.

	Depth int // Number of plies to search. 0 is ignored.
	Nodes int // Number of nodes to search. 0 is ignored.
}

// Go sends a "go" command. It starts engine calculations.
func (c *Client) Go(p GoParameters) {
	fmt.Fprintf(c.w, "go")
	if p.Ponder {
		fmt.Fprintf(c.w, " ponder")
	}
	if p.Infinite {
		fmt.Fprintf(c.w, " infinite")
	}
	if p.Mate > 0 {
		fmt.Fprintf(c.w, " mate %d", p.Mate)
	}
	if p.MoveTime > 0 {
		fmt.Fprintf(c.w, " movetime %d", p.MoveTime.Milliseconds())
	}
	if p.WhiteTime > 0 {
		fmt.Fprintf(c.w, " wtime %d", p.WhiteTime.Milliseconds())
	}
	if p.BlackTime > 0 {
		fmt.Fprintf(c.w, " btime %d", p.BlackTime.Milliseconds())
	}
	if p.WhiteIncrement > 0 {
		fmt.Fprintf(c.w, " winc %d", p.WhiteIncrement.Milliseconds())
	}
	if p.BlackIncrement > 0 {
		fmt.Fprintf(c.w, " binc %d", p.BlackIncrement.Milliseconds())
	}
	if p.MovesToGo > 0 {
		fmt.Fprintf(c.w, " movestogo %d", p.MovesToGo)
	}
	if p.Depth > 0 {
		fmt.Fprintf(c.w, " depth %d", p.Depth)
	}
	if p.Nodes > 0 {
		fmt.Fprintf(c.w, " nodes %d", p.Nodes)
	}
	// For best compatibility, "searchmoves" is in the final position.
	if len(p.SearchMoves) > 0 {
		fmt.Fprintf(c.w, " searchmoves %s", strings.Join(p.SearchMoves, " "))
	}
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
