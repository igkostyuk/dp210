package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	ticketLength = 6

	// ErrCountMethod indicates that program called with unknown count method.
	ErrCountMethod = errors.New("unknown count method")
)

func countDigits(n int) int {
	count := 0
	for n != 0 {
		n, count = n/10, count+1
	}
	return count
}

func isMoskowLucky(n int) bool {
	l, r, dc := 0, 0, countDigits(n)
	for i := 0; i < dc/2; i++ {
		l += n / int(math.Pow10(dc-1-i)) % 10
		r += n / int(math.Pow10(i)) % 10
	}
	return l == r
}

func isPiterLucky(n int) bool {
	o, e, d, dc := 0, 0, 0, countDigits(n)
	for i := 0; i < dc; i++ {
		d = n / int(math.Pow10(i)) % 10
		if d%2 == 0 {
			e += d
			continue
		}
		o += d
	}
	return o == e
}

func countNumbers(min, max int, countMethod string) (int, error) {
	switch countMethod {
	case "Moskow":
		return countLucky(min, max, isMoskowLucky), nil
	case "Piter":
		return countLucky(min, max, isPiterLucky), nil
	default:
		return 0, ErrCountMethod
	}
}

func countLucky(min, max int, isLucky func(int) bool) int {
	luckyCounter := 0
	for i := min; i <= max; i++ {
		if isLucky(i) {
			luckyCounter++
		}
	}
	return luckyCounter
}

func readTicketNumber(r io.Reader) (int, error) {
	line, err := readLine(r)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(line)
	if err != nil || n < 0 || len(line) != ticketLength {
		return 0, strconv.ErrSyntax
	}
	return n, nil
}

func readCountingMethod(r io.Reader) (string, error) {
	line, err := readLine(r)
	if err != nil {
		return "", fmt.Errorf("read filename: %w", err)
	}
	data, err := os.ReadFile(line)
	if err != nil {
		return "", fmt.Errorf("read config: %w", err)
	}
	return strings.TrimSpace(string(data)), err
}

func readLine(r io.Reader) (string, error) {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\n"), nil
}

// Task count lucky tickets from to number from reader
// with method in file.
func Task(r io.Reader, w io.Writer) error {
	fmt.Fprint(w, "Enter config filename: ")
	method, err := readCountingMethod(r)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "count method: %s\n", method)

	fmt.Fprint(w, "Min: ")
	min, err := readTicketNumber(r)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "Max: ")
	max, err := readTicketNumber(r)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "--Result--")
	c, err := countNumbers(min, max, method)
	if err != nil {
		return err
	}
	fmt.Fprint(w, c)
	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: counting lucky numbers\n", os.Args[0])
}

func main() {
	usage(os.Stdout)
	if err := Task(os.Stdin, os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
