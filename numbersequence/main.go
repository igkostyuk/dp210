package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

var ErrNumberSyntax = errors.New("number should be positive float")

func usage() {
	fmt.Fprintf(os.Stdout, "%s: print numeric sequence till square number\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <number>", os.Args[0])
}

func WriteSequence(w io.Writer, n float64) {
	n = math.Sqrt(n)
	number := int(n)
	if n-math.Trunc(n) > 0 {
		number++
	}
	if number != 0 {
		fmt.Fprint(w, 0)
	}
	for i := 1; i < number; i++ {
		fmt.Fprint(w, ",", i)
	}
}

func Task() error {
	n, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil || n < 0 {
		return ErrNumberSyntax
	}
	WriteSequence(os.Stdout, n)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(0)
	}
	if err := Task(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
