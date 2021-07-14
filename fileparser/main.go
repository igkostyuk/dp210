package main

import (
	"bytes"
	"fmt"
	"os"
)

func replace(filename, oldLine, newLine string) error {
	src, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("replace read file: %w", err)
	}
	src = bytes.ReplaceAll(src, []byte(oldLine), []byte(newLine))
	return os.WriteFile(filename, src, 0o644)
}

func count(filename, line string) (int, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return 0, fmt.Errorf("count read file: %w", err)
	}

	return bytes.Count(src, []byte(line)), nil
}

func printLinesCount(filename, line string) error {
	linesCount, err := count(filename, line)
	if err != nil {
		return err
	}
	fmt.Printf("in the file:%s, the line:\"%s\" appears %d times", filename, line, linesCount)

	return nil
}

func usage() {
	fmt.Fprintf(os.Stdout, "%s: count line or replace line in file\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <filename> <line>", os.Args[0])
	fmt.Fprintf(os.Stdout, "usage: %s <filename> <oldline> <newline>", os.Args[0])
}

func Task() error {
	switch len(os.Args) {
	case 3:
		return printLinesCount(os.Args[1], os.Args[2])
	case 4:
		return replace(os.Args[1], os.Args[2], os.Args[3])
	default:
		usage()

		return nil
	}
}

func main() {
	if err := Task(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
