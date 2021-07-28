package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/igkostyuk/dp210/envelopes/envelope"
)

type size struct {
	name  string
	value float64
}

type sizePair [2]size

type envelopsPair [2]*envelope.Envelope

func getSizePairValues(r *bufio.Reader, w io.Writer, sp *sizePair) error {
	for i, s := range sp {
		fmt.Fprintf(w, "%s: ", s.name)
		text, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading input:%w", err)
		}
		text = strings.TrimSuffix(text, "\n")
		sp[i].value, err = strconv.ParseFloat(text, 64)
		if err != nil {
			return fmt.Errorf("parsing size:%w", err)
		}
	}
	return nil
}

func getEnvelops(r *bufio.Reader, w io.Writer, sps []sizePair) ([]*envelope.Envelope, error) {
	envs := make([]*envelope.Envelope, 0)
	for _, sp := range sps {
		name := sp[0].name + sp[1].name
		fmt.Fprintf(w, "Enter %s envelope sizes\n", name)
		if err := getSizePairValues(r, w, &sp); err != nil {
			return nil, err
		}
		env, err := envelope.NewEnvelope(name, sp[0].value, sp[1].value)
		if err != nil {
			return nil, fmt.Errorf("get envelopes:%w", err)
		}
		envs = append(envs, env)

	}
	return envs, nil
}

func getFitEnvelops(r *bufio.Reader, w io.Writer, sps []sizePair) ([]*envelopsPair, error) {
	envs, err := getEnvelops(r, w, sps)
	if err != nil {
		return nil, fmt.Errorf("get fit envelopes:%w", err)
	}
	fp := make([]*envelopsPair, 0)
	for i := range envs {
		for j := range envs {
			if i != j && envs[i].IsFitsIn(envs[j]) {
				fp = append(fp, &envelopsPair{envs[i], envs[j]})
			}
		}
	}
	return fp, nil
}

func confirm(r *bufio.Reader, confirms []string) bool {
	text, err := r.ReadString('\n')
	if err == nil {
		text = strings.TrimSuffix(text, "\n")
		for _, confirm := range confirms {
			if strings.EqualFold(text, confirm) {
				return true
			}
		}
	}
	return false
}

// Task check if one of readed envelops can fit in other.
func Task(r io.Reader, w io.Writer) {
	br := bufio.NewReader(r)
	sps := []sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}}
	done, confirms := false, []string{"y", "yes"}
	for !done {
		eps, err := getFitEnvelops(br, w, sps)
		switch {
		case err != nil:
			fmt.Fprintln(w, err)
		case len(eps) == 0:
			fmt.Fprintln(w, "envelops can't fit")
		default:
			for _, ep := range eps {
				fmt.Fprintf(w, "envelope %s can fit in %s\n", ep[0], ep[1])
			}
		}
		fmt.Fprintf(w, "continue %v ?:", confirms)
		done = !confirm(br, confirms)
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: checks if one envelope can fit in another\n", os.Args[0])
}

func main() {
	usage(os.Stdout)
	Task(os.Stdin, os.Stdout)
}
