package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

var (
	// ErrNumberSyntax indicates that a value does not have the right syntax for the size type.
	ErrNumberSyntax = errors.New("number should be positive float")
	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("parameter length should be 1 <number>")
)

// WriteSequence write natural numbers separated by commas,
// the square of which is less than a given n.
func WriteSequence(w io.Writer, n float64) error {
	bw := bufio.NewWriter(w)
	s := math.Sqrt(n)
	number := int(s)
	if s-math.Trunc(s) > 0 {
		number++
	}
	if number != 0 {
		fmt.Fprint(bw, 0)
	}
	for i := 1; i < number; i++ {
		fmt.Fprint(bw, ",", i)
	}
	return bw.Flush()
}

// Task write natural numbers separated by commas,
// the square of which is less than main param.
func Task(w io.Writer, args []string) error {
	if len(args) != 1 {
		return ErrParameters
	}
	n, err := strconv.ParseFloat(args[0], 64)
	if err != nil || n < 0 {
		return ErrNumberSyntax
	}
	return WriteSequence(w, n)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print numeric sequence till square number\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <number>", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
		}
		fmt.Println(err)
	}
}
