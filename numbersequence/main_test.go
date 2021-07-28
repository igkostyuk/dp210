package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteSequence(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{"10", args{10}, "0,1,2,3", assert.NoError},
		{"25", args{25}, "0,1,2,3,4", assert.NoError},
		{"negative number", args{-25}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, WriteSequence(w, tt.args.n))
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
		{"empty params", args{[]string{}}, "", assert.Error},
		{"valid int params", args{[]string{"4"}}, "0,1", assert.NoError},
		{"valid float params", args{[]string{"4.1"}}, "0,1,2", assert.NoError},
		{"invalid float params", args{[]string{"invalid"}}, "", assert.Error},
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
	name := "test"
	os.Args[0] = name
	tests := []struct {
		name  string
		wantW string
	}{
		{"usage",
			fmt.Sprintf("%s: print numeric sequence till square number\n", name) +
				fmt.Sprintf("usage: %s <number>", name)},
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
