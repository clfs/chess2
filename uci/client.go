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

	Result Result // The result of the last search. Set using Client.Go.
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
func (c *Client) UCI() (name, author string, opts []Option, err error) {
	fmt.Fprintln(c.w, "uci")

	s := bufio.NewScanner(c.r)
outer:
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "id name "):
			name = strings.TrimPrefix(line, "id name ")
		case strings.HasPrefix(line, "id author "):
			author = strings.TrimPrefix(line, "id author ")
		case strings.HasPrefix(line, "option "):
			var opt Option
			if err := opt.UnmarshalText([]byte(line)); err != nil {
				return "", "", nil, err
			}
			opts = append(opts, opt)
		case line == "uciok":
			break outer
		}
	}

	return name, author, opts, s.Err()
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

// Search contains parameters for the "go" command. Note that fields of type
// time.Duration are truncated to the millisecond.
type Search struct {
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

func (s Search) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "go")
	if s.Ponder {
		fmt.Fprintf(&b, " ponder")
	}
	if s.Infinite {
		fmt.Fprintf(&b, " infinite")
	}
	if s.Mate != 0 {
		fmt.Fprintf(&b, " mate %d", s.Mate)
	}
	if s.MoveTime != 0 {
		fmt.Fprintf(&b, " movetime %d", s.MoveTime)
	}
	if s.WhiteTime != 0 {
		fmt.Fprintf(&b, " wtime %d", s.WhiteTime)
	}
	if s.BlackTime != 0 {
		fmt.Fprintf(&b, " btime %d", s.BlackTime)
	}
	if s.WhiteIncrement != 0 {
		fmt.Fprintf(&b, " winc %d", s.WhiteIncrement)
	}
	if s.BlackIncrement != 0 {
		fmt.Fprintf(&b, " binc %d", s.BlackIncrement)
	}
	if s.MovesToGo != 0 {
		fmt.Fprintf(&b, " movestogo %d", s.MovesToGo)
	}
	if s.Depth != 0 {
		fmt.Fprintf(&b, " depth %d", s.Depth)
	}
	if s.Nodes != 0 {
		fmt.Fprintf(&b, " nodes %d", s.Nodes)
	}
	// For best compatibility, "searchmoves" is in the final position.
	if len(s.SearchMoves) > 0 {
		fmt.Fprintf(&b, "searchmoves %s", strings.Join(s.SearchMoves, " "))
	}
	return b.String()
}

// Go sends a "go" command. It starts engine calculations.
func (c *Client) Go(s Search) {
	fmt.Fprintf(c.w, "%s\n", s)
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
