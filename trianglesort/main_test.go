package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/igkostyuk/dp210/trianglesort/triangle"
	"github.com/stretchr/testify/assert"
)

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
		{"with spaces", args{" t e   s t "}, "test"},
		{"with uppercase", args{"tEsT"}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sanitize(tt.args.text))
		})
	}
}

func Test_parseTriangle(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name      string
		args      args
		want      *triangle.Triangle
		assertion assert.ErrorAssertionFunc
	}{

		{
			"valid triangle",
			args{"test,1,1,1"},
			&triangle.Triangle{Name: "test", A: 1, B: 1, C: 1},
			assert.NoError,
		},
		{
			"negative triangle size",
			args{"test,1,-1,1"},
			nil, assert.Error,
		},
		{
			"invalid triangle size",
			args{"test,1,invalid,1"},
			nil, assert.Error,
		},
		{
			"invalid triangle params too long",
			args{"test,1,1,1,1"},
			nil, assert.Error,
		},
		{
			"invalid triangle params too short",
			args{"test,1,1"},
			nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTriangle(tt.args.text)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getTriangle(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name      string
		args      args
		want      *triangle.Triangle
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid triangle",
			args{r: bufio.NewReader(strings.NewReader("test,1,1,1\n"))},
			&triangle.Triangle{Name: "test", A: 1, B: 1, C: 1}, assert.NoError,
		},
		{
			"invalid reader",
			args{r: bufio.NewReader(iotest.ErrReader(errors.New("error")))},
			nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTriangle(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isConfirm(t *testing.T) {
	type args struct {
		text     string
		confirms []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"confirm",
			args{"test", []string{"test"}}, true,
		},
		{
			"not confirm",
			args{"test", []string{"no"}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isConfirm(tt.args.text, tt.args.confirms))
		})
	}
}

func TestConfirm(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		wantW string
	}{
		{
			"confirm",
			args{bufio.NewReader(strings.NewReader("yes\n"))}, true, "continue [y yes] ?:",
		},
		{
			"not confirm",
			args{bufio.NewReader(strings.NewReader("test\n"))}, false, "continue [y yes] ?:",
		},
		{
			"error reader",
			args{bufio.NewReader(iotest.ErrReader(errors.New("test")))}, false, "continue [y yes] ?:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			assert.Equal(t, tt.want, Confirm(tt.args.r, w))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{

		{"invalid triangle", args{strings.NewReader("test,INVALID,1,1\n no\n")},
			"Enter triangle <name>,<a>,<b>,<c>: " +
				"parsing triangle size:size should be positive float\n" +
				"continue [y yes] ?:", assert.NoError},

		{"valid triangle", args{strings.NewReader("test,1,1,1\n no\n")},
			"Enter triangle <name>,<a>,<b>,<c>: continue [y yes] ?:" +
				"============= Triangles list: ===============\n" +
				"1. [Triangle test]: 0.43 cm", assert.NoError},

		{"valid two triangles", args{strings.NewReader("one,1,1,1\n yes\ntwo,5,4,3\n no\n")},
			"Enter triangle <name>,<a>,<b>,<c>: continue [y yes] ?:" +
				"Enter triangle <name>,<a>,<b>,<c>: continue [y yes] ?:" +
				"============= Triangles list: ===============\n" +
				"1. [Triangle two]: 6.00 cm" +
				"2. [Triangle one]: 0.43 cm", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, Task(tt.args.r, w))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_usage(t *testing.T) {
	name := "test"
	os.Args[0] = name

	tests := []struct {
		name  string
		wantW string
	}{
		{"usage", fmt.Sprintf("%s: sorts the triangles from user input\n", name)},
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
		{"no panic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
