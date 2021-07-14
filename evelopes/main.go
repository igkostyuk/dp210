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
	Confirms      = []string{"y", "yes"}
	AnswerCDFit   = "CD envelope can fit into AB envelope"
	AnswerABFit   = "AB envelope can fit into CD envelope"
	AnswerCantFit = "can't fit"

	ErrSizeSyntax = errors.New("size should be positive float")
)

type Envelope struct {
	a float64
	b float64
}

func NewEnvelope(a, b float64) *Envelope {
	return &Envelope{a, b}
}

func (r *Envelope) isFits(br *Envelope) bool {
	a, b, q, p := br.a, br.b, r.a, r.b
	// Rectangles in Rectangles John E. Wetzel
	// https://www.researchgate.net/publication/273573138_Rectangles_in_Rectangles(page 204)
	if b > a {
		a, b = b, a
	}
	if q > p {
		q, p = p, q
	}

	return q < b && (p < a || b*(p*p+q*q) > (2*p*q*a+(p*p-q*q)*math.Sqrt(p*p+q*q-a*a)))
}

func scanFloatSize(s *bufio.Scanner, message string) (float64, error) {
	fmt.Print(message)
	if s.Scan() {
		return strconv.ParseFloat(s.Text(), 64)
	}

	return 0, s.Err()
}

func scanSizes(s *bufio.Scanner) (a, b, c, d float64, err error) {
	fmt.Println("Enter AB envelope and CD  envelope  sizes")
	if a, err = scanFloatSize(s, "A: "); err != nil || a <= 0 {
		err = fmt.Errorf("size A: %w", ErrSizeSyntax)

		return
	}
	if b, err = scanFloatSize(s, "B: "); err != nil || b <= 0 {
		err = fmt.Errorf("size B: %w", ErrSizeSyntax)

		return
	}
	if c, err = scanFloatSize(s, "C: "); err != nil || c <= 0 {
		err = fmt.Errorf("size C: %w", ErrSizeSyntax)

		return
	}
	if d, err = scanFloatSize(s, "D: "); err != nil || d <= 0 {
		err = fmt.Errorf("size A: %w", ErrSizeSyntax)

		return
	}
	err = s.Err()

	return
}

func isDone(s *bufio.Scanner, confirms []string) bool {
	fmt.Printf("continue %v ?:", confirms)
	if !s.Scan() {
		return true
	}
	for _, word := range confirms {
		if strings.EqualFold(s.Text(), word) {
			return false
		}
	}

	return true
}

func Task(s *bufio.Scanner) error {
	var done bool
	for !done {
		a, b, c, d, err := scanSizes(s)
		if err != nil {
			return err
		}
		ab := NewEnvelope(a, b)
		cd := NewEnvelope(c, d)
		switch {
		case cd.isFits(ab):
			fmt.Println(AnswerCDFit)
		case ab.isFits(cd):
			fmt.Println(AnswerABFit)
		default:
			fmt.Println(AnswerCantFit)
		}

		done = isDone(s, Confirms)
	}

	return s.Err()
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: checks if one envelope can fit in another\n", os.Args[0])
}

func main() {
	usage()
	scanner := bufio.NewScanner(os.Stdin)
	if err := Task(scanner); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
