package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	confirms             = []string{"y", "yes"}
	triangleParamsString = "<name>,<a>,<b>,<c>"

	ErrSizeSyntax       = errors.New("size should be positive float")
	ErrParametersLength = errors.New("expected 4 comma separated parameters")
	ErrParametersValid  = errors.New("not valid sizes: expected a + b <= c or a + c <= b or b + c <= a")
)

type Triangle struct {
	name    string
	a, b, c float64
	area    float64
}

func NewTriangle(name string, a, b, c float64) *Triangle {
	// Heron's formula
	s := (a + b + c) / 2
	area := math.Sqrt(s * (s - a) * (s - b) * (s - c))

	return &Triangle{name: name, a: a, b: b, c: c, area: area}
}

func sanitize(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")

	return strings.ToLower(text)
}

func getTriangle(s *bufio.Scanner, msg string) (*Triangle, error) {
	fmt.Printf("Enter triangle %s: ", msg)
	if s.Scan() {
		params := strings.Split(sanitize(s.Text()), ",")
		name, a, b, c, err := parseParams(params)
		if err != nil {
			return nil, err
		}

		if a+b <= c || a+c <= b || b+c <= a {
			return nil, ErrParametersValid
		}

		return NewTriangle(name, a, b, c), nil
	}

	return nil, s.Err()
}

func parseParams(params []string) (name string, a, b, c float64, err error) {
	if len(params) != 4 {
		err = ErrParametersLength

		return
	}
	name = params[0]
	a, err = strconv.ParseFloat(params[1], 64)
	if err != nil || a <= 0 {
		err = fmt.Errorf("a :%w", ErrSizeSyntax)

		return
	}

	b, err = strconv.ParseFloat(params[2], 64)
	if err != nil || b <= 0 {
		err = fmt.Errorf("b :%w", ErrSizeSyntax)

		return
	}

	c, err = strconv.ParseFloat(params[3], 64)
	if err != nil || c <= 0 {
		err = fmt.Errorf("c :%w", ErrSizeSyntax)
	}

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

func printTriangles(ts []*Triangle) {
	fmt.Println("============= Triangles list: ===============")
	for i, t := range ts {
		fmt.Printf("%d. [Triangle %s]: %.2f cm\n", i+1, t.name, t.area)
	}
}

func Task(s *bufio.Scanner) error {
	var ts []*Triangle
	var done bool
	for !done {
		t, err := getTriangle(s, triangleParamsString)
		if err != nil {
			fmt.Println(err)
		}
		if err == nil && t != nil {
			ts = append(ts, t)
		}
		done = isDone(s, confirms)
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].area > ts[j].area })
	printTriangles(ts)

	return s.Err()
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: sorts the triangles from user input\n", os.Args[0])
}

func main() {
	usage()
	scanner := bufio.NewScanner(os.Stdin)
	if err := Task(scanner); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
