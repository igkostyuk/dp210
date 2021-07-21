package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func Test_getTripletName(t *testing.T) {
	type args struct {
		n          int
		dictionary map[int]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"valid three word three digits number",
			args{555, numbersDictionary},
			"пятьсот пятьдесят пять", false,
		},
		{
			"valid two word three digits number",
			args{717, numbersDictionary},
			"семьсот семнадцать", false,
		},
		{
			"missing hundreds in dictionary",
			args{100, map[int]string{}},
			"", true,
		},
		{
			"missing tens in dictionary",
			args{20, map[int]string{}},
			"", true,
		},
		{
			"missing teens in dictionary",
			args{10, map[int]string{}},
			"", true,
		},
		{
			"missing ones in dictionary",
			args{21, map[int]string{20: "test"}},
			"", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTripletName(tt.args.n, tt.args.dictionary)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTripletName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTripletName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTriplets(t *testing.T) {
	type args struct {
		number int64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"one triplet number", args{number: 123}, []int{123}},
		{"three triplet number", args{number: 123_456_789}, []int{789, 456, 123}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTriplets(tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTriplets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNumberName(t *testing.T) {
	type args struct {
		number string
		nd     map[int]string
		pd     [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"invalid number",
			args{number: "invalid", nd: numbersDictionary, pd: periodDictionary},
			"", true,
		},
		{
			"number zero",
			args{number: "0", nd: numbersDictionary, pd: periodDictionary},
			"нуль", false,
		},
		{
			"negative one",
			args{number: "-1", nd: numbersDictionary, pd: periodDictionary},
			"минус один", false,
		},
		{
			"first triplet missing in dictionary",
			args{number: "10", nd: map[int]string{}, pd: periodDictionary},
			"", true,
		},
		{
			"second triplet missing in dictionary",
			args{number: "1000", nd: map[int]string{0: "test"}, pd: periodDictionary},
			"", true,
		},
		{
			"\"ин\"suffix",
			args{number: "1000", nd: numbersDictionary, pd: periodDictionary},
			"одна тысяча", false,
		},
		{
			"\"ва\"suffix",
			args{number: "2000", nd: numbersDictionary, pd: periodDictionary},
			"две тысячи", false,
		},
		{
			"second plural index number",
			args{number: "11000", nd: numbersDictionary, pd: periodDictionary},
			"одиннадцать тысяч", false,
		},

		{
			"three triplet number",
			args{number: "123456789", nd: numbersDictionary, pd: periodDictionary},
			"сто двадцать три миллиона четыреста пятьдесят шесть тысяч семьсот восемьдесят девять", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNumberName(tt.args.number, tt.args.nd, tt.args.pd)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNumberName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNumberName() = %v, want %v", got, tt.want)
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
		{"invalid args", args{[]string{"invalid"}}, "", true},
		{"invalid args lenght", args{[]string{"2", "2"}}, "", true},
		{"valid args", args{[]string{"2"}}, "2 - два\n", false},
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
	}{
		{"usage", "test: print number name\n" + "usage: test <number>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			usage(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("usage() got\n%v\nwant\n%v", gotW, tt.wantW)
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
