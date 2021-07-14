package main

import (
	"errors"
	"fmt"
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

func (br *Board) String() string {
	var b strings.Builder
	for _, r := range br.squares {
		b.WriteString(string(r))
		b.WriteRune('\n')
	}

	return b.String()
}

func NewBoard(height, width int, black, white rune) *Board {
	squares := make([][]rune, height)
	var c, cc, n, nc rune
	c, n = black, white
	for i := range squares {
		c, cc, n, nc = n, n, c, c
		squares[i] = make([]rune, width)
		for j := range squares[i] {
			cc, nc = nc, cc
			squares[i][j] = cc
		}
	}

	return &Board{height: height, width: width, squares: squares}
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: print chessboard\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <height> <width>", os.Args[0])
}

func Task(args []string) error {
	height, err := strconv.Atoi(args[1])
	if err != nil || height <= 0 {
		return fmt.Errorf("height: %w", ErrSize)
	}
	width, err := strconv.Atoi(args[2])
	if err != nil || width <= 0 {
		return fmt.Errorf("width: %w", ErrSize)
	}
	board := NewBoard(height, width, blackSymbol, whiteSymbol)
	fmt.Print(board)

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
