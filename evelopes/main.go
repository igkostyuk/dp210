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
	// Confirms words of confirmation.
	Confirms       = []string{"y", "yes"}
	positiveAnswer = "%s envelope can fit into %s envelope\n"
	negativeAnswer = "can't fit\n"
	// ErrSizeSyntax indicates that a value does not have the right syntax for the size type.
	ErrSizeSyntax = errors.New("size should be positive float")
)

// Envelope represent envelope.
type Envelope struct {
	name   string
	height float64
	width  float64
}

// NewEnvelope create envelope with name height width sizes.
func NewEnvelope(name string, height, width float64) (*Envelope, error) {
	if height <= 0 || width <= 0 {
		return nil, ErrSizeSyntax
	}
	return &Envelope{name: name, height: height, width: width}, nil
}

func (r *Envelope) isFits(br *Envelope) bool {
	a, b := br.width, br.height
	q, p := r.width, r.height
	// Rectangles in Rectangles John E. Wetzel
	// https://www.researchgate.net/publication/
	// 273573138_Rectangles_in_Rectangles(page 204)
	if b > a {
		a, b = b, a
	}
	if q > p {
		q, p = p, q
	}
	if q > b {
		return false
	}
	if p < a {
		return true
	}
	// fits diagonally
	return b*(p*p+q*q) > (2*p*q*a + (p*p-q*q)*math.Sqrt(p*p+q*q-a*a))
}

func getEnvelope(br *bufio.Reader, w io.Writer, pair [2]string) (*Envelope, error) {
	name, vals := pair[0]+pair[1], [2]float64{}
	fmt.Fprintf(w, "Enter %s envelope sizes\n", name)
	for j, name := range pair {
		fmt.Fprintf(w, "%s: ", name)

		text, err := br.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("reading envelope input:%w", err)
		}

		vals[j], err = strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			return nil, fmt.Errorf("parsing size : %w", err)
		}
	}
	return NewEnvelope(name, vals[0], vals[1])
}

func getEnvelops(r io.Reader, w io.Writer, sizeNames [][2]string) ([]*Envelope, error) {
	envelops := make([]*Envelope, len(sizeNames))
	var err error
	br := bufio.NewReader(r)
	for i, pair := range sizeNames {
		if envelops[i], err = getEnvelope(br, w, pair); err != nil {
			return nil, fmt.Errorf("get envelopes:%w", err)
		}
	}
	return envelops, nil
}

// Confirm write question
// and return true if reads one of the words of confirmation.
func Confirm(r io.Reader, w io.Writer, confirms []string) bool {
	fmt.Fprintf(w, "continue %v ?:", confirms)
	br := bufio.NewReader(r)
	text, err := br.ReadString('\n')
	if err == nil {
		text = strings.TrimSpace(text)
		for _, word := range confirms {
			if strings.EqualFold(text, word) {
				return true
			}
		}
	}
	return false
}

// Task check if one of readed envelops can fit in other.
func Task(r io.Reader, w io.Writer) {
	sizeNames := [2][2]string{{"A", "B"}, {"C", "D"}}
	var done bool
	for !done {
		envs, err := getEnvelops(r, w, sizeNames[:])
		switch {
		case err != nil:
			fmt.Fprintln(w, err)
		case envs[0].isFits(envs[1]):
			fmt.Fprintf(w, positiveAnswer, envs[0].name, envs[1].name)
		case envs[1].isFits(envs[0]):
			fmt.Fprintf(w, positiveAnswer, envs[1].name, envs[0].name)
		default:
			fmt.Fprint(w, negativeAnswer)
		}
		done = !Confirm(r, w, Confirms)
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: checks if one envelope can fit in another\n", os.Args[0])
}

func main() {
	usage(os.Stdout)
	Task(os.Stdin, os.Stdout)
}
