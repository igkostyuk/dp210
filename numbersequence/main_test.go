package main

import (
	"bytes"
	"testing"
)

func TestWriteSequence(t *testing.T) {
	type args struct {
		n float64
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{"10", args{10}, "0,1,2,3"},
		{"25", args{25}, "0,1,2,3,4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WriteSequence(w, tt.args.n); err != nil {
				t.Errorf("Task() error = %v", err)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteSequence() = %v, want %v", gotW, tt.wantW)
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
		{"empty params", args{[]string{}}, "", true},
		{"valid int params", args{[]string{"4"}}, "0,1", false},
		{"valid float params", args{[]string{"4.1"}}, "0,1,2", false},
		{"invalid float params", args{[]string{"invalid"}}, "", true},
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
	wantW := ""
	w := &bytes.Buffer{}
	usage(w)
	if gotW := w.String(); gotW != wantW {
		t.Errorf("usage() = %v, want %v", gotW, wantW)
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
