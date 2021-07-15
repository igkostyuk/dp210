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
	blackSymbol = '*'
	whiteSymbol = ' '

	ErrSize = errors.New("size should be a positive integer")
)

type Board struct {
	height  int
	width   int
	squares [][]rune
}

func NewBoard(height, width int, squares [][]rune) *Board {
	return &Board{height: height, width: width, squares: squares}
}

func (br *Board) String() string {
	var b strings.Builder
	for _, r := range br.squares {
		b.WriteString(string(r))
		b.WriteRune('\n')
	}

	return b.String()
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: print chessboard\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <height> <width>", os.Args[0])
}

func WriteBoard(w io.Writer, height, width int, blackSymbol, whiteSymbol rune) {
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

	board := NewBoard(height, width, squares)
	fmt.Fprint(w, board)
}

func ParseParams(args []string) (height, width int, err error) {
	height, err = strconv.Atoi(args[1])
	if err != nil || height <= 0 {
		err = fmt.Errorf("height: %w", ErrSize)

		return
	}
	width, err = strconv.Atoi(args[2])
	if err != nil || width <= 0 {
		err = fmt.Errorf("width: %w", ErrSize)

		return
	}

	return height, width, nil
}

func Task(args []string) error {
	height, width, err := ParseParams(args)
	if err != nil {
		return fmt.Errorf("parsing parameters:%w", err)
	}
	WriteBoard(os.Stdout, height, width, blackSymbol, whiteSymbol)

	return nil
}

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(0)
	}
	if err := Task(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
