package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/igkostyuk/dp210/trianglesort/triangle"
)

var (
	// ErrParametersLength indicates that program called with wrong number of parameters
	ErrParametersLength = errors.New("expected 4 comma separated parameters")
)

func sanitize(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")
	return strings.ToLower(text)
}

func parseTriangle(text string) (*triangle.Triangle, error) {
	var err error
	params := strings.Split(sanitize(text), ",")
	if len(params) != 4 {
		return nil, ErrParametersLength
	}
	var size [3]float64
	for i, p := range params[1:] {
		if size[i], err = strconv.ParseFloat(p, 64); err != nil {
			return nil, fmt.Errorf("parsing triangle size:%w", triangle.ErrSizeSyntax)
		}
	}
	return triangle.NewTriangle(params[0], size[0], size[1], size[2])
}
func getTriangle(r *bufio.Reader) (*triangle.Triangle, error) {
	text, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("get triangle reading string:%w", err)
	}
	text = strings.TrimSuffix(text, "\n")
	return parseTriangle(text)
}

func isConfirm(text string, confirms []string) bool {
	for _, confirm := range confirms {
		if strings.EqualFold(text, confirm) {
			return true
		}
	}
	return false
}

// Confirm write question
// and return true if reads one of the words of confirmation.
func Confirm(r *bufio.Reader, w io.Writer) bool {
	confirms := []string{"y", "yes"}
	fmt.Fprintf(w, "continue %v ?:", confirms)
	text, err := r.ReadString('\n')
	if err == nil {
		return isConfirm(strings.TrimSpace(text), confirms)
	}
	return false
}

// Task read tringles from strign and write sorter by area.
func Task(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	var (
		ts   []*triangle.Triangle
		done bool
	)
	for !done {
		fmt.Fprintf(w, "Enter triangle <name>,<a>,<b>,<c>: ")
		switch t, err := getTriangle(br); err {
		case nil:
			ts = append(ts, t)
		default:
			fmt.Fprintln(w, err)
		}
		done = !Confirm(br, w)
	}
	if len(ts) > 0 {
		sort.Slice(ts, func(i, j int) bool { return ts[i].Area() > ts[j].Area() })
		fmt.Fprintln(w, "============= Triangles list: ===============")
	}
	for i, t := range ts {
		fmt.Fprintf(w, "%d. %s\n", i+1, t)
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
	}
}
