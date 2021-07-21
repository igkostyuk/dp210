package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestNewEnvelope(t *testing.T) {
	type args struct {
		name   string
		height float64
		width  float64
	}
	tests := []struct {
		name    string
		args    args
		want    *Envelope
		wantErr bool
	}{
		{"valid params", args{"AB", 1.1, 2.2}, &Envelope{name: "AB", height: 1.1, width: 2.2}, false},
		{"negative first param", args{"AB", -1.1, 2.2}, nil, true},
		{"negative second param", args{"AB", 1.1, -2.2}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnvelope(tt.args.name, tt.args.height, tt.args.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEnvelope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnvelope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvelope_isFits(t *testing.T) {
	type sizes struct {
		height float64
		width  float64
	}
	tests := []struct {
		name string
		ab   sizes
		cd   sizes
		want bool
	}{
		{
			"size bigger than argument size",
			sizes{10, 10}, sizes{5, 5}, false,
		},
		{
			"size smaller than argument size",
			sizes{5, 5}, sizes{10, 10}, true,
		},
		{
			"height bigger than argument size",
			sizes{9, 4}, sizes{5, 10}, true,
		},
		{
			"width bigger than argument size",
			sizes{4, 9}, sizes{10, 5}, true,
		},
		{
			"diagonal fit",
			sizes{10, 1}, sizes{9, 9}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := &Envelope{
				name:   "ab",
				height: tt.ab.height,
				width:  tt.ab.width,
			}
			cd := &Envelope{
				name:   "cd",
				height: tt.cd.height,
				width:  tt.cd.width,
			}

			if got := ab.isFits(cd); got != tt.want {
				t.Errorf("Envelope.isFits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvelope(t *testing.T) {
	type args struct {
		r         *bufio.Reader
		sizeNames [2]string
	}
	tests := []struct {
		name    string
		args    args
		want    *Envelope
		wantW   string
		wantErr bool
	}{
		{
			"valid int sizes",
			args{bufio.NewReader(strings.NewReader("1" + "\n" + "1" + "\n")), [2]string{"A", "B"}},
			&Envelope{"AB", 1, 1}, "Enter AB envelope sizes\nA: B: ", false,
		},
		{
			"valid int sizes first with spaces",
			args{bufio.NewReader(strings.NewReader(" 1 " + "\n" + "1" + "\n")), [2]string{"A", "B"}},
			&Envelope{"AB", 1, 1}, "Enter AB envelope sizes\nA: B: ", false,
		},
		{
			"valid float sizes",
			args{bufio.NewReader(strings.NewReader("1.1" + "\n" + "1.1" + "\n")), [2]string{"C", "D"}},
			&Envelope{"CD", 1.1, 1.1}, "Enter CD envelope sizes\nC: D: ", false,
		},
		{
			"invalid first float sizes",
			args{bufio.NewReader(strings.NewReader("invalid" + "\n" + "1" + "\n")), [2]string{"A", "B"}},
			nil, "Enter AB envelope sizes\nA: ", true,
		},
		{
			"invalid second float sizes",
			args{bufio.NewReader(strings.NewReader("1" + "\n" + "invalid" + "\n")), [2]string{"A", "B"}},
			nil, "Enter AB envelope sizes\nA: B: ", true,
		},
		{
			"negative first float sizes",
			args{bufio.NewReader(strings.NewReader("-1" + "\n" + "1" + "\n")), [2]string{"A", "B"}},
			nil, "Enter AB envelope sizes\nA: B: ", true,
		},
		{
			"negative second float sizes",
			args{bufio.NewReader(strings.NewReader("1" + "\n" + "-1" + "\n")), [2]string{"A", "B"}},
			nil, "Enter AB envelope sizes\nA: B: ", true,
		},
		{
			"reader error",
			args{bufio.NewReader(iotest.ErrReader(errors.New("error"))), [2]string{"A", "B"}},
			nil, "Enter AB envelope sizes\nA: ", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := getEnvelope(tt.args.r, w, tt.args.sizeNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEnvelope() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvelope() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("getEnvelope() = \n%v\n, want \n%v", gotW, tt.wantW)
			}
		})
	}
}

func Test_Confirm(t *testing.T) {
	type args struct {
		r        io.Reader
		confirms []string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		wantW string
	}{
		{
			"confirmed",
			args{strings.NewReader("y\n"), []string{"y"}},
			true, "continue [y] ?:"},
		{
			"not confirmed",
			args{strings.NewReader("n\n"), []string{"y"}},
			false, "continue [y] ?:"},
		{
			"confirmed in upper register",
			args{strings.NewReader("Y\n"), []string{"y"}},
			true, "continue [y] ?:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := Confirm(tt.args.r, w, tt.args.confirms); got != tt.want {
				t.Errorf("isDone() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("isDone() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			"AB fit",
			args{
				strings.NewReader("1\n" + "1\n" + "2\n" + "2\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"AB envelope can fit into CD envelope\n" +
				"continue [y yes] ?:",
		},
		{
			"CD fit",
			args{
				strings.NewReader("2\n" + "2\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"CD envelope can fit into AB envelope\n" +
				"continue [y yes] ?:",
		},
		{
			"cant fit",
			args{
				strings.NewReader("1\n" + "1\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"can't fit\n" +
				"continue [y yes] ?:",
		},
		{
			"invalid first envelope param",
			args{
				strings.NewReader("INVALID\n" + "1\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"get envelopes:parsing size : strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:",
		},
		{
			"invalid second param",
			args{
				strings.NewReader("1\n" + "INVALID\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"get envelopes:parsing size : strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:",
		},
		{
			"invalid third param",
			args{
				strings.NewReader("1" + "\n" + "1" + "\n" + "INVALID" + "\n" + "1" + "\n" + "no" + "\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"get envelopes:parsing size : strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			Task(tt.args.r, w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() got\n%v\nwant\n%v", gotW, tt.wantW)
			}
		})
	}
}

func Test_usage(t *testing.T) {
	tests := []struct {
		name  string
		wantW string
	}{
		{"usage", "test: checks if one envelope can fit in another\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args[0] = "test"
			w := &bytes.Buffer{}
			usage(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() got\n%v\nwant\n%v", gotW, tt.wantW)
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
