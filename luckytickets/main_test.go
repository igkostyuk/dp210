package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

func Test_isMoskowLucky(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"101", args{101}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMoskowLucky(tt.args.n); got != tt.want {
				t.Errorf("isMoskowLucky() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPiterLucky(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"112", args{112}, true},
		{"211", args{211}, true},
		{"121211", args{121211}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPiterLucky(tt.args.n); got != tt.want {
				t.Errorf("isPiterLucky() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countNumbers(t *testing.T) {
	type args struct {
		min         int
		max         int
		countMethod string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"moscow", args{min: 120123, max: 320320, countMethod: "Moskow"}, 11187, false},
		{"piter", args{min: 120123, max: 320320, countMethod: "Piter"}, 5790, false},
		{"invalid", args{min: 120123, max: 320320, countMethod: "invalid"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countNumbers(tt.args.min, tt.args.max, tt.args.countMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("countNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readTicketNumber(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid number", args{strings.NewReader("000001\n")}, 1, false},
		{"too short number", args{strings.NewReader("00001\n")}, 0, true},
		{"too long number", args{strings.NewReader("1000001\n")}, 0, true},
		{"invalid", args{strings.NewReader("invalid\n")}, 0, true},
		{"err reader", args{iotest.ErrReader(errors.New("test"))}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTicketNumber(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTicketNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readTicketNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readCountingMethod(t *testing.T) {

	f, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	// nolint: errcheck
	defer os.Remove(f.Name()) // clean up
	if _, err := f.WriteString("test method"); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"temp file", args{strings.NewReader(f.Name() + "\n")}, "test method", false},
		{"reader error", args{iotest.ErrReader(errors.New("error"))}, "", true},
		{"invalid file", args{strings.NewReader("invalid\n")}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCountingMethod(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCountingMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readCountingMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask(t *testing.T) {

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Task(tt.args.r, w); (err != nil) != tt.wantErr {
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
	tests := []struct {
		name  string
		wantW string
	}{
		{"usage", fmt.Sprintf("%s: counting lucky numbers\n", os.Args[0])},
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
