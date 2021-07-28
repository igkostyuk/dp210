package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replace(t *testing.T) {

	repeats := 12
	tempFileName := createTempFileWithText(strings.Repeat("test", repeats))
	defer func() { // clean up
		err := os.Remove(tempFileName)
		if err != nil {
			log.Fatal(err)
		}
	}()

	type args struct {
		filename string
		oldLine  string
		newLine  string
	}
	tests := []struct {
		name      string
		args      args
		want      int
		assertion assert.ErrorAssertionFunc
	}{
		{"invalid filename", args{"INVALID", "old", "new"}, 0, assert.Error},
		{"temp file", args{tempFileName, "test", "replacer"}, repeats, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replace(tt.args.filename, tt.args.oldLine, tt.args.newLine)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_count(t *testing.T) {

	repeats := 20
	tempFileName := createTempFileWithText(strings.Repeat("test", repeats))
	defer func() { // clean up
		err := os.Remove(tempFileName)
		if err != nil {
			log.Fatal(err)
		}
	}()

	type args struct {
		filename string
		line     string
	}
	tests := []struct {
		name      string
		args      args
		want      int
		assertion assert.ErrorAssertionFunc
	}{
		{"invalid filename", args{"INVALID", "line"}, 0, assert.Error},
		{"temp file", args{tempFileName, "test"}, repeats, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := count(tt.args.filename, tt.args.line)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTask(t *testing.T) {

	os.Args[0] = "test"

	repeats := 10
	word, newWord := "test", "count"
	tempFileName := createTempFileWithText(strings.Repeat(word, repeats))
	defer func() { // clean up
		err := os.Remove(tempFileName)
		if err != nil {
			log.Fatal(err)
		}
	}()

	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"invalid file count",
			args{[]string{"invalid", "line"}},
			"", assert.Error,
		},
		{
			"invalid file replace",
			args{[]string{"invalid", "old line", "new line"}},
			"", assert.Error,
		},
		{
			"temp file count",
			args{[]string{tempFileName, word}},
			fmt.Sprintf("in the file:%s,", tempFileName) +
				fmt.Sprintf(" the line:\"%s\" appears %d times", word, repeats),
			assert.NoError,
		},
		{
			"temp file replace",
			args{[]string{tempFileName, word, newWord}},
			fmt.Sprintf("in the file:%s,", tempFileName) +
				fmt.Sprintf(" the line:\"%s\" replaced", word) +
				fmt.Sprintf(" with line:\"%s\" %d times", newWord, repeats),
			assert.NoError,
		},
		{
			"usage", args{[]string{}},
			"test: count line or replace line in file\n" +
				"usage: test <filename> <line>\n" +
				"usage: test <filename> <oldline> <newline>\n",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, Task(w, tt.args.args))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_usage(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
	}{
		{
			"usage",
			"test: count line or replace line in file\n" +
				"usage: test <filename> <line>\n" +
				"usage: test <filename> <oldline> <newline>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			usage(w)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"main don't panic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func createTempFileWithText(text string) string {
	f, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(text); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}
