package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	ticketLength = 6

	ErrNegativeNumber = errors.New("number is negative")
	ErrTicketLength   = errors.New("unexpected ticket length")
	ErrCountMethod    = errors.New("unknown count method")
)

func getTicketNumber(number string, tl int) (int, error) {
	n, err := strconv.Atoi(number)
	if err != nil {
		return 0, strconv.ErrSyntax
	}
	if len(number) != tl {
		return 0, fmt.Errorf("number is't %d digit long:%w", tl, ErrTicketLength)
	}

	if n < 0 {
		return 0, ErrNegativeNumber
	}

	return n, nil
}

func isMoskowLucky(n int) bool {
	l, r := 0, 0
	for i := 0; i < ticketLength/2; i++ {
		l += n / int(math.Pow10(ticketLength-1-i)) % 10
		r += n / int(math.Pow10(i)) % 10
	}

	return l == r
}

func isPiterLucky(n int) bool {
	var o, e, d int
	for i := 0; i < ticketLength; i++ {
		d = n / int(math.Pow10(i)) % 10
		if d%2 == 0 {
			e += d

			continue
		}
		o += d
	}

	return o == e
}

func countLuckyNumbers(min, max int, isLucky func(int) bool) int {
	luckyCounter := 0
	for i := min; i <= max; i++ {
		if isLucky(i) {
			luckyCounter++
		}
	}

	return luckyCounter
}

func getCountingMethod(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("read config: %w", err)
	}

	return strings.TrimSpace(string(data)), err
}

func getCountFunc(countMethod string) (func(int) bool, error) {
	switch countMethod {
	case "Moskow":
		return isMoskowLucky, nil
	case "Piter":
		return isPiterLucky, nil
	default:
		return nil, ErrCountMethod
	}
}

func Task() (err error) {
	var min, max int
	var countMethod string
	var countFunc func(int) bool
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter config filename: ")
	if scanner.Scan() {
		if countMethod, err = getCountingMethod(scanner.Text()); err != nil {
			return err
		}
		if countFunc, err = getCountFunc(countMethod); err != nil {
			return err
		}
		fmt.Printf("count method: %s\n", countMethod)
	}

	fmt.Print("Min: ")
	if scanner.Scan() {
		if min, err = getTicketNumber(scanner.Text(), ticketLength); err != nil {
			return err
		}
	}

	fmt.Print("Max: ")
	if scanner.Scan() {
		if max, err = getTicketNumber(scanner.Text(), ticketLength); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("task scanner:%w", err)
	}
	fmt.Printf("\n--Result--\n %d\n", countLuckyNumbers(min, max, countFunc))

	return nil
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: counting lucky numbers\n", os.Args[0])
}

func main() {
	usage()
	if err := Task(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
