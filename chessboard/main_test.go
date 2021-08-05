package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseParameters(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		want      *Parameters
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid parameters", args{[]string{"1", "2"}},
			&Parameters{Height: 1, Width: 2}, assert.NoError,
		},
		{
			"too many parameters", args{[]string{"1", "2", "3"}},
			nil, assert.Error,
		},
		{
			"invalid height parameter", args{[]string{"INVALID", "2"}},
			nil, assert.Error,
		},
		{
			"invalid width parameter", args{[]string{"1", "INVALID"}},
			nil, assert.Error,
		},
		{
			"negative height parameter", args{[]string{"-1", "2"}},
			nil, assert.Error,
		},
		{
			"negative width parameter", args{[]string{"1", "-2"}},
			nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseParameters(tt.args.args)
			tt.assertion(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{"valid params", args{[]string{"1", "1"}}, "*\n", assert.NoError},
		{"invalid params", args{[]string{"invalid", "1"}}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, run(w, tt.args.args))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		p *Parameters
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid arguments",
			args{&Parameters{Height: 3, Width: 3}},
			"* *\n * \n* *\n", assert.NoError,
		},
		{
			"negative arguments",
			args{&Parameters{Height: -3, Width: 3}},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, Task(w, tt.args.p))
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
		{
			"test name",
			fmt.Sprintf("%s: print chessboard\nusage: %s <height> <width>\n", name, name),
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
		{"run main"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
