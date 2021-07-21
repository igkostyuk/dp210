package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_replace(t *testing.T) {

	repeats := 12
	tmpfile, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	// nolint: errcheck
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.WriteString(strings.Repeat("test", repeats)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	type args struct {
		filename string
		oldLine  string
		newLine  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"invalid filename", args{"INVALID", "old", "new"}, 0, true},
		{"temp file", args{tmpfile.Name(), "test", "replacer"}, repeats, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replace(tt.args.filename, tt.args.oldLine, tt.args.newLine)
			if (err != nil) != tt.wantErr {
				t.Errorf("replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_count(t *testing.T) {

	repeats := 12
	tmpfile, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	// nolint: errcheck
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.WriteString(strings.Repeat("test", repeats)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	type args struct {
		filename string
		line     string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"invalid filename", args{"INVALID", "line"}, 0, true},
		{"temp file", args{tmpfile.Name(), "test"}, repeats, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := count(tt.args.filename, tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask(t *testing.T) {

	repeats := 12
	tmpfile, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	// nolint: errcheck
	defer os.Remove(tmpfile.Name()) // clean up
	if _, err := tmpfile.WriteString(strings.Repeat("test", repeats)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	os.Args[0] = "test"

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"invalid file count", args{[]string{"invalid", "line"}}, "", true},
		{"invalid file replace", args{[]string{"invalid", "old line", "new line"}}, "", true},
		{"temp file count", args{[]string{tmpfile.Name(), "test"}},
			"in the file:" + tmpfile.Name() + ", the line:\"test\"" +
				" appears " + "12" + " times", false},
		{"temp file replace", args{[]string{tmpfile.Name(), "test", "count"}},
			"in the file:" + tmpfile.Name() +
				", the line:\"test\" replaced" +
				" with line:\"count\" " + "12" + " times", false},

		{"usage", args{[]string{}},
			"test: count line or replace line in file\n" +
				"usage: test <filename> <line>\n" +
				"usage: test <filename> <oldline> <newline>\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Task(w, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Task() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() =\n%v\n, want \n%v\n", gotW, tt.wantW)
			}
		})
	}
}

func Test_usage(t *testing.T) {

	os.Args[0] = "test"

	tests := []struct {
		name  string
		wantW string
	}{{
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
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("usage() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
