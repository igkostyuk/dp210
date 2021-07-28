package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appendFromDictionary(t *testing.T) {
	type args struct {
		key   int
		words *[]string
		d     NumberDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      *[]string
		assertion assert.ErrorAssertionFunc
	}{
		{"number in dictionary", args{10, &[]string{}, nd}, &[]string{nd[10]}, assert.NoError},
		{"number not in dictionary", args{55, &[]string{}, nd}, &[]string{}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, appendFromDictionary(tt.args.key, tt.args.words, tt.args.d))
			assert.Equal(t, tt.want, tt.args.words)
		})
	}
}

func Test_getTripletName(t *testing.T) {
	type args struct {
		n int
		d NumberDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid three word three digits number",
			args{555, nd},
			"пятьсот пятьдесят пять", assert.NoError,
		},
		{
			"valid two word three digits number",
			args{717, nd},
			"семьсот семнадцать", assert.NoError,
		},
		{
			"missing hundreds in dictionary",
			args{100, map[int]string{}},
			"", assert.Error,
		},
		{
			"missing tens in dictionary",
			args{20, map[int]string{}},
			"", assert.Error,
		},
		{
			"missing teens in dictionary",
			args{10, map[int]string{}},
			"", assert.Error,
		},
		{
			"missing ones in dictionary",
			args{21, map[int]string{20: "test"}},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTripletName(tt.args.n, tt.args.d)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
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
			assert.Equal(t, tt.want, parseTriplets(tt.args.number))
		})
	}
}

func Test_getPeriodName(t *testing.T) {
	type args struct {
		idx int
		n   int
		pd  PeriodDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{"valid name", args{1, 1, pd}, "миллион", assert.NoError},
		{"too big period index", args{6, 1, pd}, "", assert.Error},
		{"too big period index", args{0, 5, PeriodDictionary{{""}}}, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPeriodName(tt.args.idx, tt.args.n, tt.args.pd)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_appendTripletsNames(t *testing.T) {
	type args struct {
		tr    []int
		words *[]string
		nd    NumberDictionary
		pd    PeriodDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      *[]string
		assertion assert.ErrorAssertionFunc
	}{

		{
			"three triplets",
			args{[]int{789, 456, 123}, &[]string{}, nd, pd},
			&[]string{
				"сто двадцать три", "миллиона",
				"четыреста пятьдесят шесть", "тысяч",
				"семьсот восемьдесят девять",
			}, assert.NoError,
		},
		{
			"\"ин\"suffix",
			args{[]int{0, 1}, &[]string{}, nd, pd},
			&[]string{"одна", "тысяча", ""}, assert.NoError,
		},
		{
			"\"ва\"suffix",
			args{[]int{0, 2}, &[]string{}, nd, pd},
			&[]string{"две", "тысячи", ""}, assert.NoError,
		},
		{
			"first triplet missing in dictionary",
			args{[]int{1}, &[]string{}, nil, pd},
			&[]string{}, assert.Error,
		},
		{
			"second triplet missing in dictionary",
			args{[]int{1, 2}, &[]string{}, NumberDictionary{1: "one"}, pd},
			&[]string{}, assert.Error,
		},
		{
			"missing period name",
			args{[]int{1, 2}, &[]string{}, NumberDictionary{1: "one", 2: "two"}, nil},
			&[]string{}, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, appendTripletsNames(tt.args.tr, tt.args.words, tt.args.nd, tt.args.pd))
			assert.Equal(t, tt.want, tt.args.words)
		})
	}
}

func Test_getNumberName(t *testing.T) {
	type args struct {
		n  int64
		nd NumberDictionary
		pd PeriodDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"number zero",
			args{n: 0, nd: nd, pd: pd},
			"нуль", assert.NoError,
		},
		{
			"negative one",
			args{n: -1, nd: nd, pd: pd},
			"минус один", assert.NoError,
		},
		{
			"get triplet name error",
			args{n: 1, nd: nil, pd: pd},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNumberName(tt.args.n, tt.args.nd, tt.args.pd)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name       string
		args       args
		wantW      string
		dictionary NumberDictionary
		assertion  assert.ErrorAssertionFunc
	}{
		{"invalid args", args{[]string{"invalid"}}, "", nd, assert.Error},
		{"invalid args length", args{[]string{"2", "2"}}, "", nd, assert.Error},
		{"valid args", args{[]string{"2"}}, "2 - два\n", nd, assert.NoError},
		{"invalid dictionary", args{[]string{"2"}}, "", nil, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nd = tt.dictionary
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
		{"usage test name", fmt.Sprintf("%s: print number name\nusage: %s <number>", name, name)},
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
