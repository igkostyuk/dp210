package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func replace(filename, oldLine, newLine string) (int, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("replace read file: %w", err)
	}

	count := bytes.Count(src, []byte(oldLine))
	src = bytes.ReplaceAll(src, []byte(oldLine), []byte(newLine))

	return count, os.WriteFile(filename, src, 0o644)
}

func count(filename, line string) (int, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("count read file: %w", err)
	}

	return bytes.Count(src, []byte(line)), nil
}

// Task count or replace text in file depends of params.
func Task(w io.Writer, args []string) error {
	switch len(args) {
	case 2:
		filename, line := args[0], args[1]
		count, err := count(filename, line)
		if err != nil {
			return fmt.Errorf("tast count: %w", err)
		}
		fmt.Fprintf(w, "in the file:%s,", filename)
		fmt.Fprintf(w, " the line:\"%s\" appears %d times", line, count)
		return nil
	case 3:
		filename, oldLine, newLine := args[0], args[1], args[2]
		count, err := replace(filename, oldLine, newLine)
		if err != nil {
			return fmt.Errorf("task replace: %w", err)
		}
		fmt.Fprintf(w, "in the file:%s,", filename)
		fmt.Fprintf(w, " the line:\"%s\" replaced", oldLine)
		fmt.Fprintf(w, " with line:\"%s\" %d times", newLine, count)
		return nil
	default:
		usage(w)
		return nil
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "%s: count line or replace line in file\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <filename> <line>\n", os.Args[0])
	fmt.Fprintf(w, "usage: %s <filename> <oldline> <newline>\n", os.Args[0])
}

func main() {
	if err := Task(os.Stdout, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
