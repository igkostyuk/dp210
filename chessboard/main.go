package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/igkostyuk/dp210/chessboard/board"
)

var (

	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("should be 2 parameters <height> <width>")
)

// Parameters represent task parameters.
type Parameters struct {
	Width  int
	Height int
}

func parseParameters(args []string) (*Parameters, error) {
	if len(args) != 2 {
		return nil, ErrParameters
	}

	height, err := strconv.Atoi(args[0])
	if err != nil || height <= 0 {
		return nil, fmt.Errorf("parse param height: %w", board.ErrSize)
	}

	width, err := strconv.Atoi(args[1])
	if err != nil || width <= 0 {
		return nil, fmt.Errorf("parse param width: %w", board.ErrSize)
	}

	return &Parameters{Width: width, Height: height}, nil
}

func run(w io.Writer, args []string) error {
	p, err := parseParameters(args)
	if err != nil {
		return fmt.Errorf("parsing parameters:%w", err)
	}
	return Task(w, p)
}

// Task write string of board with task parameters.
func Task(w io.Writer, p *Parameters) error {
	b, err := board.NewBoard(p.Height, p.Width, board.BlackSymbol, board.WhiteSymbol)
	if err != nil {
		return fmt.Errorf("task creating board:%w", err)
	}
	return b.Write(w)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print chessboard\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <height> <width>\n", os.Args[0])
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
		}
		fmt.Println(err)
	}
}
