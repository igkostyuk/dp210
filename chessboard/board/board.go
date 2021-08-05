package board

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	// BlackSymbol represent black square symbol
	BlackSymbol = '*'
	// WhiteSymbol represent white square symbol
	WhiteSymbol = ' '

	// ErrSize indicates that a value does not have the right syntax for the size type.
	ErrSize = errors.New("size should be a positive integer")
)

// Board represent chess board as rune matrix.
type Board struct {
	Height  int
	Width   int
	Squares [][]rune
}

// NewBoard creates new board with height and width sizes and black and whiter symbols.
func NewBoard(height, width int, blackSymbol, whiteSymbol rune) (*Board, error) {
	if height <= 0 || width <= 0 {
		return nil, ErrSize
	}

	squares := createSquares(height, width, blackSymbol, whiteSymbol)

	return &Board{Height: height, Width: width, Squares: squares}, nil
}

// String return string representation of board.
func (br *Board) String() string {
	var b strings.Builder
	for _, r := range br.Squares {
		b.WriteString(string(r))
		b.WriteRune('\n')
	}

	return b.String()
}

// Write write board to writer
func (br *Board) Write(w io.Writer) error {
	bw := bufio.NewWriter(w)
	for j, row := range br.Squares {
		for i := range row {
			if _, err := bw.WriteRune(br.Squares[j][i]); err != nil {
				return fmt.Errorf("writing board:%w", err)
			}
		}
		if _, err := bw.WriteRune('\n'); err != nil {
			return fmt.Errorf("writing board:%w", err)
		}
	}
	return bw.Flush()
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
