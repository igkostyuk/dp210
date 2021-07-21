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

func TestNewTriangle(t *testing.T) {
	type args struct {
		name string
		a    float64
		b    float64
		c    float64
	}
	tests := []struct {
		name    string
		args    args
		want    *Triangle
		wantErr bool
	}{
		{
			"valid triangle",
			args{name: "test", a: 1, b: 1, c: 1}, &Triangle{name: "test", a: 1, b: 1, c: 1}, false,
		},
		{
			"negative size",
			args{name: "test", a: -1, b: 1, c: 1}, nil, true,
		},
		{
			"too big size",
			args{name: "test", a: 1, b: 5, c: 1}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTriangle(tt.args.name, tt.args.a, tt.args.b, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTriangle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTriangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTriangle_Area(t *testing.T) {
	type fields struct {
		a float64
		b float64
		c float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"same sizes #1", fields{a: 5, b: 4, c: 3}, 6},
		{"same sizes #2", fields{a: 4, b: 5, c: 3}, 6},
		{"same sizes #3", fields{a: 3, b: 4, c: 5}, 6},
		{"all sizes 1", fields{a: 1, b: 1, c: 1}, 0.4330127018922193},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Triangle{
				a: tt.fields.a,
				b: tt.fields.b,
				c: tt.fields.c,
			}
			if got := tr.Area(); got != tt.want {
				t.Errorf("Triangle.Area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTriangle_String(t *testing.T) {
	type fields struct {
		name string
		a    float64
		b    float64
		c    float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"round area", fields{name: "test", a: 3, b: 4, c: 5}, "[Triangle test]: 6.00 cm"},
		{"float area", fields{name: "test", a: 1, b: 1, c: 1}, "[Triangle test]: 0.43 cm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Triangle{
				name: tt.fields.name,
				a:    tt.fields.a,
				b:    tt.fields.b,
				c:    tt.fields.c,
			}
			if got := tr.String(); got != tt.want {
				t.Errorf("Triangle.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *reader
	}{
		{"stdin", args{os.Stdin}, &reader{br: bufio.NewReader(os.Stdin)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newReader(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sanitize(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with tabs", args{"\tt\tes\tt"}, "test"},
		{"with newline", args{"\nt\nest\n"}, "test"},
		{"with spaces", args{" t e   s t "}, "test"},
		{"with uppercase", args{"tEsT"}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitize(tt.args.text); got != tt.want {
				t.Errorf("sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_ReadTriangle(t *testing.T) {
	type fields struct {
		br *bufio.Reader
	}
	type args struct {
		delim byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Triangle
		wantErr bool
	}{
		{
			"valid triangle",
			fields{br: bufio.NewReader(strings.NewReader("test,1,1,1\n"))},
			args{'\n'}, &Triangle{name: "test", a: 1, b: 1, c: 1}, false,
		},
		{
			"negative triangle size",
			fields{br: bufio.NewReader(strings.NewReader("test,1,-1,1\n"))},
			args{'\n'}, nil, true,
		},
		{
			"invalid triangle size",
			fields{br: bufio.NewReader(strings.NewReader("test,1,invalid,1\n"))},
			args{'\n'}, nil, true,
		},
		{
			"invalid triangle params too long",
			fields{br: bufio.NewReader(strings.NewReader("test,1,1,1,1\n"))},
			args{'\n'}, nil, true,
		},
		{
			"invalid triangle params too short",
			fields{br: bufio.NewReader(strings.NewReader("test,1,1\n"))},
			args{'\n'}, nil, true,
		},
		{
			"error reader",
			fields{br: bufio.NewReader(iotest.ErrReader(errors.New("test error")))},
			args{'\n'}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				br: tt.fields.br,
			}
			got, err := r.ReadTriangle(tt.args.delim)
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.ReadTriangle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reader.ReadTriangle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_confirm(t *testing.T) {
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
			if got := confirm(tt.args.r, w, tt.args.confirms); got != tt.want {
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
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"invalid triangle", args{strings.NewReader("test,INVALID,1,1\n yes\n yes\n")},
			"Enter triangle <name>,<a>,<b>,<c>: " +
				"size should be positive float\n" +
				"continue [y yes] ?:" +
				"============= Triangles list: ===============\n", false},

		{"valid triangle", args{strings.NewReader("test,1,1,1\n yes\n yes\n")},
			"Enter triangle <name>,<a>,<b>,<c>: continue [y yes] ?:" +
				"============= Triangles list: ===============\n" +
				"1. [Triangle test]: 0.43 cm", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Task(tt.args.r, w); (err != nil) != tt.wantErr {
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
	}{
		{"usage", "test: sorts the triangles from user input\n"},
	}
	for _, tt := range tests {
		w := &bytes.Buffer{}
		t.Run(tt.name, func(t *testing.T) {
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
