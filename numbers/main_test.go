package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertTensToWords(t *testing.T) {
	type args struct {
		t int
		d NumberDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"one word number",
			args{1, nd}, []string{"один"}, assert.NoError,
		},
		{
			"two words number",
			args{21, nd}, []string{"двадцать", "один"}, assert.NoError,
		},
		{
			"not in dictionary one word number",
			args{1, map[int]string{}}, nil, assert.Error,
		},
		{
			"not in dictionary first word of two word number",
			args{21, map[int]string{1: "one"}}, nil, assert.Error,
		},
		{
			"not in dictionary second word of two word number",
			args{21, map[int]string{20: "twenty"}}, nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTensToWords(tt.args.t, tt.args.d)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_convertTripletToWords(t *testing.T) {
	type args struct {
		t int
		d NumberDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid three words number",
			args{123, nd}, []string{"сто", "двадцать", "три"}, assert.NoError,
		},
		{
			"not in dictionary first of  three words number",
			args{123, map[int]string{20: "", 3: ""}}, nil, assert.Error,
		},
		{
			"not in dictionary third of  three words number",
			args{123, map[int]string{100: "", 20: ""}}, nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTripletToWords(tt.args.t, tt.args.d)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_convertTripletToWord(t *testing.T) {
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
			"valid three words number",
			args{123, nd}, "сто двадцать три", assert.NoError,
		},
		{
			"empty dictionary",
			args{123, map[int]string{}}, "", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTripletToWord(tt.args.n, tt.args.d)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
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
		{
			"valid name",
			args{1, 10, pd}, "миллионов", assert.NoError,
		},
		{
			"valid name",
			args{1, 3, pd}, "миллиона", assert.NoError,
		},
		{
			"valid name",
			args{1, 1, pd}, "миллион", assert.NoError,
		},
		{
			"valid name",
			args{1, 50, pd}, "миллионов", assert.NoError,
		},
		{
			"too big period index",
			args{6, 1, pd}, "", assert.Error,
		},
		{
			"too big index in period",
			args{1, 5, PeriodDictionary{{""}}}, "", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPeriodName(tt.args.idx, tt.args.n, tt.args.pd)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_fixThousandsSuffix(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"ин suffix", "один", "одна"},
		{"ин in words", "одиннадцать", "одиннадцать"},
		{"ва suffix", "два", "две"},
		{"ва in words", "двадцать", "двадцать"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixThousandsSuffix(&tt.args)
			assert.Equal(t, tt.want, tt.args)
		})
	}
}

func Test_convertTripletsToWords(t *testing.T) {
	type args struct {
		tr []int
		nd NumberDictionary
		pd PeriodDictionary
	}
	tests := []struct {
		name      string
		args      args
		want      []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"three triplets",
			args{[]int{789, 456, 123}, nd, pd},
			[]string{
				"сто двадцать три", "миллиона",
				"четыреста пятьдесят шесть", "тысяч",
				"семьсот восемьдесят девять",
			}, assert.NoError,
		},
		{
			"\"ин\"suffix",
			args{[]int{0, 1}, nd, pd},
			[]string{"одна", "тысяча"}, assert.NoError,
		},
		{
			"\"ва\"suffix",
			args{[]int{0, 2}, nd, pd},
			[]string{"две", "тысячи"}, assert.NoError,
		},
		{
			"zero triplets in middle",
			args{[]int{0, 0, 1}, nd, pd},
			[]string{"один", "миллион"}, assert.NoError,
		},
		{
			"first triplet missing in dictionary",
			args{[]int{1}, nil, pd},
			nil, assert.Error,
		},
		{
			"second triplet missing in dictionary",
			args{[]int{1, 2}, NumberDictionary{1: "one"}, pd},
			nil, assert.Error,
		},
		{
			"missing period name",
			args{[]int{1, 2}, NumberDictionary{1: "one", 2: "two"}, nil},
			nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTripletsToWords(tt.args.tr, tt.args.nd, tt.args.pd)
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

func Test_convertNumberToWord(t *testing.T) {
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
			"invalid number dictionary",
			args{n: 1, nd: nil, pd: pd},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertNumberToWord(tt.args.n, tt.args.nd, tt.args.pd)
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
		{"valid args", args{[]string{"2"}}, "2 - два\n", nd, assert.NoError},
		{"invalid args number", args{[]string{"invalid"}}, "", nd, assert.Error},
		{"invalid args length", args{[]string{"2", "2"}}, "", nd, assert.Error},
		{"invalid dictionary", args{[]string{"2"}}, "", nil, assert.Error},
	}
	tmp := NumberDictionary{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dictionary == nil {
				tmp, nd = nd, tmp
			}
			w := &bytes.Buffer{}
			tt.assertion(t, Task(w, tt.args.args))
			assert.Equal(t, tt.wantW, w.String())
			if tt.dictionary == nil {
				tmp, nd = nd, tmp
			}

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
		{"usage test name", fmt.Sprintf("%s: print number converted to words\nusage: %s <number>\n", name, name)},
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
