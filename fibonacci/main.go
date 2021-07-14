package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var ErrNumberSyntax = errors.New("number should be positive int")

func usage() {
	fmt.Fprintf(os.Stdout, "%s: print fibonacci in the specified range\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <number> <number>", os.Args[0])
}

func fibonacciSeq() func() int {
	a, b := 0, 1

	return func() int {
		res := a
		a, b = b, a+b

		return res
	}
}

func WriteFibonacciSequence(w io.Writer, from, to int) {
	n := 0
	if from > to {
		from, to = to, from
	}
	nextInt := fibonacciSeq()
	printed := false
	for n < to {
		if n >= from {
			if printed {
				fmt.Fprint(w, ",", n)
			}
			if !printed {
				printed = !printed
				fmt.Fprint(w, n)
			}
		}
		n = nextInt()
	}
}

func Task() (err error) {
	var fn, sn int
	if fn, err = strconv.Atoi(os.Args[1]); err != nil || fn < 0 {
		return ErrNumberSyntax
	}
	if sn, err = strconv.Atoi(os.Args[2]); err != nil || sn < 0 {
		return ErrNumberSyntax
	}
	WriteFibonacciSequence(os.Stdout, fn, sn)

	return nil
}

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(0)
	}
	if err := Task(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
