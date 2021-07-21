package main

import (
	"bytes"
	"os"
	"testing"
)

func TestWriteFibonacciSequence(t *testing.T) {
	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"first 19", args{0, 2600}, "0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584", false},
		{"first reverse params 19", args{2600, 0}, "0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584", false},
		{"first 4-19", args{2, 2600}, "2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WriteFibonacciSequence(w, tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("WriteFibonacciSequence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteFibonacciSequence() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"valid parameters", args{[]string{"0", "2"}}, "0,1,1", false},
		{"invalid first parameters", args{[]string{"invalid", "2"}}, "", true},
		{"invalid second parameters", args{[]string{"2", "invalid"}}, "", true},
		{"invalid parameters lenght", args{[]string{"2", "3", "4"}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Task(w, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Task() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_usage(t *testing.T) {

	os.Args[0] = "test"

	tests := []struct {
		name  string
		wantW string
	}{
		{"usage", "test: print fibonacci in the specified range\n" + "usage: test <number> <number>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			usage(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() = %v, want %v", gotW, tt.wantW)
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
