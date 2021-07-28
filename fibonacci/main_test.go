package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteFibonacciSequence(t *testing.T) {
	type args struct {
		from int
		to   int
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"first 19",
			args{0, 2600}, "0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584",
			assert.NoError,
		},
		{
			"first reverse params 19",
			args{2600, 0}, "0,1,1,2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584",
			assert.NoError,
		},
		{
			"first 4-19",
			args{2, 2600}, "2,3,5,8,13,21,34,55,89,144,233,377,610,987,1597,2584",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, WriteFibonacciSequence(w, tt.args.from, tt.args.to))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func TestTask(t *testing.T) {
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
			"valid parameters",
			args{[]string{"0", "2"}}, "0,1,1", assert.NoError,
		},
		{
			"invalid first parameters",
			args{[]string{"invalid", "2"}}, "", assert.Error,
		},
		{
			"invalid second parameters",
			args{[]string{"2", "invalid"}}, "", assert.Error,
		},
		{
			"invalid parameters length",
			args{[]string{"2", "3", "4"}}, "", assert.Error,
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
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"main no panic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
