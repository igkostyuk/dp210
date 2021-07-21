package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	// ErrNumberSyntax indicates that a value does not have the right
	// syntax for the parameters number type.
	ErrNumberSyntax = errors.New("number should be positive int")
	// ErrParameters indicates that program called with wrong number of parameters
	ErrParameters = errors.New("parameter length should be 2 <number> <number>")
)

func fibonacciSeq() func() int {
	a, b := 0, 1

	return func() int {
		res := a
		a, b = b, a+b

		return res
	}
}

// WriteFibonacciSequence write fibonacci sequence
// from from number till to number.
func WriteFibonacciSequence(w io.Writer, from, to int) error {
	bw := bufio.NewWriter(w)
	if from > to {
		from, to = to, from
	}
	nextInt := fibonacciSeq()
	n := nextInt()
	printed := false
	for n < to {
		if n >= from {
			if printed {
				fmt.Fprint(bw, ",", n)
			}
			if !printed {
				printed = !printed
				fmt.Fprint(bw, n)
			}
		}
		n = nextInt()
	}
	return bw.Flush()
}

// Task write fibonacci sequence from and till args params.
func Task(w io.Writer, args []string) error {
	if len(args) != 2 {
		return ErrParameters
	}
	var fn, sn int
	var err error
	if fn, err = strconv.Atoi(args[0]); err != nil || fn < 0 {
		return ErrNumberSyntax
	}
	if sn, err = strconv.Atoi(args[1]); err != nil || sn < 0 {
		return ErrNumberSyntax
	}
	return WriteFibonacciSequence(w, fn, sn)
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: print fibonacci in the specified range\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <number> <number>", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		if errors.Is(err, ErrParameters) {
			usage(os.Stdout)
			os.Exit(0)
		}
		fmt.Println(err)
		os.Exit(0)
	}
}
