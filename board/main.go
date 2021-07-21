package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	// BlackSymbol represent black square symbol
	BlackSymbol = '*'
	// WhiteSymbol represent white square symbol
	WhiteSymbol = ' '

	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("parameter length should be 2 <height> <width>")
	// ErrSize indicates that a value does not have the right syntax for the size type.
	ErrSize = errors.New("size should be a positive integer")
)

// Board represent chess board as rune matrix.
type Board struct {
	height  int
	width   int
	squares [][]rune
}

// NewBoard creates new board with height and width sizes and black and witer symbols.
func NewBoard(height, width int, blackSymbol, whiteSymbol rune) (*Board, error) {
	if height <= 0 || width <= 0 {
		return nil, ErrSize
	}

	squares := createSquares(height, width, blackSymbol, whiteSymbol)

	return &Board{height: height, width: width, squares: squares}, nil
}

// String return string representation of board.
func (br *Board) String() string {
	var b strings.Builder
	for _, r := range br.squares {
		b.WriteString(string(r))
		b.WriteRune('\n')
	}

	return b.String()
}

func createSquares(height, width int, blackSymbol, whiteSymbol rune) [][]rune {
	squares := make([][]rune, height)
	var c, cc, n, nc rune
	c, n = blackSymbol, whiteSymbol
	for i := range squares {
		c, cc, n, nc = n, n, c, c
		squares[i] = make([]rune, width)
		for j := range squares[i] {
			cc, nc = nc, cc
			squares[i][j] = cc
		}
	}

	return squares
}

func createBoard(parameters []string, blackSymbol, whiteSymbol rune) (*Board, error) {
	if len(parameters) != 2 {
		return nil, ErrParameters
	}

	height, err := strconv.Atoi(parameters[0])
	if err != nil || height <= 0 {
		return nil, fmt.Errorf("height: %w", ErrSize)
	}

	width, err := strconv.Atoi(parameters[1])
	if err != nil || width <= 0 {
		return nil, fmt.Errorf("width: %w", ErrSize)
	}

	return NewBoard(height, width, blackSymbol, whiteSymbol)
}

// Task write string of board with args params.
func Task(w io.Writer, args []string) error {
	board, err := createBoard(args, BlackSymbol, WhiteSymbol)
	if err != nil {
		return fmt.Errorf("creating board:%w", err)
	}
	fmt.Fprint(w, board)
	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print chessboard\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <height> <width>", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
		}
		fmt.Println(err)
		os.Exit(0)
	}
}
