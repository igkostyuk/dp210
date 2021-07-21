package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	confirms = []string{"y", "yes"}

	// ErrSizeSyntax indicates that a value does not have the right syntax for the size type.
	ErrSizeSyntax = errors.New("size should be positive float")
	// ErrParametersLength indicates that program called with wrong number of parameters
	ErrParametersLength = errors.New("expected 4 comma separated parameters")
	// ErrSizeRatio indicates that triangle size ratio is wrong.
	ErrSizeRatio = errors.New("not valid sizes: expected a + b <= c or a + c <= b or b + c <= a")
)

// Triangle represent triangle.
type Triangle struct {
	name    string
	a, b, c float64
}

// NewTriangle create triangle with name and sizes.
func NewTriangle(name string, a, b, c float64) (*Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, ErrSizeSyntax
	}
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, ErrSizeRatio
	}
	return &Triangle{name: name, a: a, b: b, c: c}, nil
}

// Area return triangle area.
func (t *Triangle) Area() float64 {
	// Heron's formula
	s := (t.a + t.b + t.c) / 2
	return math.Sqrt(s * (s - t.a) * (s - t.b) * (s - t.c))
}

func (t *Triangle) String() string {
	return fmt.Sprintf("[Triangle %s]: %.2f cm", t.name, t.Area())
}

type reader struct {
	br *bufio.Reader
}

func newReader(r io.Reader) *reader {
	return &reader{br: bufio.NewReader(r)}
}

func sanitize(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "\n", "")
	return strings.ToLower(text)
}

func (r *reader) ReadTriangle(delim byte) (*Triangle, error) {
	text, err := r.br.ReadString(delim)
	if err != nil {
		return nil, fmt.Errorf("reading triangle:%w", err)
	}
	params := strings.Split(sanitize(text), ",")
	if len(params) != 4 {
		return nil, ErrParametersLength
	}
	var s [3]float64
	for i, p := range params[1:] {
		if s[i], err = strconv.ParseFloat(p, 64); err != nil {
			return nil, ErrSizeSyntax
		}
	}
	return NewTriangle(params[0], s[0], s[1], s[2])

}

func confirm(r io.Reader, w io.Writer, confirms []string) bool {
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

// Task read tringles from strign and write sorter by area.
func Task(r io.Reader, w io.Writer) error {
	var (
		ts   []*Triangle
		done bool
	)
	tr := newReader(r)
	for !done {
		fmt.Fprintf(w, "Enter triangle <name>,<a>,<b>,<c>: ")
		switch t, err := tr.ReadTriangle('\n'); err {
		case nil:
			ts = append(ts, t)
		default:
			fmt.Fprintln(w, err)
		}
		done = !confirm(r, w, confirms)
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Area() > ts[j].Area() })
	fmt.Fprintln(w, "============= Triangles list: ===============")
	for i, t := range ts {
		fmt.Fprintf(w, "%d. %s", i+1, t)
	}

	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: sorts the triangles from user input\n", os.Args[0])
}

func main() {
	usage(os.Stdout)
	if err := Task(os.Stdin, os.Stdout); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
