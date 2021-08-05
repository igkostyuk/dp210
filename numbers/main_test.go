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
			name:      "one word number",
			args:      args{1, nd},
			want:      []string{"один"},
			assertion: assert.NoError,
		},
		{
			name:      "two words number",
			args:      args{21, nd},
			want:      []string{"двадцать", "один"},
			assertion: assert.NoError,
		},
		{
			name:      "not in dictionary one word number",
			args:      args{1, map[int]string{}},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "not in dictionary first word of two word number",
			args:      args{21, map[int]string{1: "one"}},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "not in dictionary second word of two word number",
			args:      args{21, map[int]string{20: "twenty"}},
			want:      nil,
			assertion: assert.Error,
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
			name:      "valid three words number",
			args:      args{123, nd},
			want:      []string{"сто", "двадцать", "три"},
			assertion: assert.NoError,
		},
		{
			name:      "not in dictionary first of  three words number",
			args:      args{123, map[int]string{20: "", 3: ""}},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "not in dictionary third of  three words number",
			args:      args{123, map[int]string{100: "", 20: ""}},
			want:      nil,
			assertion: assert.Error,
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
			name:      "valid three words number",
			args:      args{123, nd},
			want:      "сто двадцать три",
			assertion: assert.NoError,
		},
		{
			name:      "empty dictionary",
			args:      args{123, map[int]string{}},
			want:      "",
			assertion: assert.Error,
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
			name:      "valid name",
			args:      args{1, 10, pd},
			want:      "миллионов",
			assertion: assert.NoError,
		},
		{
			name:      "valid name",
			args:      args{1, 3, pd},
			want:      "миллиона",
			assertion: assert.NoError,
		},
		{
			name:      "valid name",
			args:      args{1, 1, pd},
			want:      "миллион",
			assertion: assert.NoError,
		},
		{
			name:      "valid name",
			args:      args{1, 50, pd},
			want:      "миллионов",
			assertion: assert.NoError,
		},
		{
			"too big period index",
			args{6, 1, pd}, "", assert.Error,
		},
		{
			name:      "too big index in period",
			args:      args{1, 5, PeriodDictionary{{""}}},
			want:      "",
			assertion: assert.Error,
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
			name: "three triplets",
			args: args{[]int{789, 456, 123}, nd, pd},
			want: []string{
				"сто двадцать три", "миллиона",
				"четыреста пятьдесят шесть", "тысяч",
				"семьсот восемьдесят девять",
			},
			assertion: assert.NoError,
		},
		{
			name:      "\"ин\"suffix",
			args:      args{[]int{0, 1}, nd, pd},
			want:      []string{"одна", "тысяча"},
			assertion: assert.NoError,
		},
		{
			name:      "\"ва\"suffix",
			args:      args{[]int{0, 2}, nd, pd},
			want:      []string{"две", "тысячи"},
			assertion: assert.NoError,
		},
		{
			name:      "zero triplets in middle",
			args:      args{[]int{0, 0, 1}, nd, pd},
			want:      []string{"один", "миллион"},
			assertion: assert.NoError,
		},
		{
			name:      "first triplet missing in dictionary",
			args:      args{[]int{1}, nil, pd},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "second triplet missing in dictionary",
			args:      args{[]int{1, 2}, NumberDictionary{1: "one"}, pd},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "missing period name",
			args:      args{[]int{1, 2}, NumberDictionary{1: "one", 2: "two"}, nil},
			want:      nil,
			assertion: assert.Error,
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
		{
			name: "one triplet number",
			args: args{number: 123},
			want: []int{123},
		},
		{
			name: "three triplet number",
			args: args{number: 123_456_789},
			want: []int{789, 456, 123}},
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
			name:      "number zero",
			args:      args{n: 0, nd: nd, pd: pd},
			want:      "нуль",
			assertion: assert.NoError,
		},
		{
			name:      "negative one",
			args:      args{n: -1, nd: nd, pd: pd},
			want:      "минус один",
			assertion: assert.NoError,
		},
		{
			name:      "invalid number dictionary",
			args:      args{n: 1, nd: nil, pd: pd},
			want:      "",
			assertion: assert.Error,
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
		{
			name:       "valid args",
			args:       args{[]string{"2"}},
			wantW:      "2 - два\n",
			dictionary: nd,
			assertion:  assert.NoError,
		},
		{
			name:       "invalid args number",
			args:       args{[]string{"invalid"}},
			wantW:      "",
			dictionary: nd,
			assertion:  assert.Error,
		},
		{
			name:       "invalid args length",
			args:       args{[]string{"2", "2"}},
			wantW:      "",
			dictionary: nd,
			assertion:  assert.Error,
		},
		{
			name:       "invalid dictionary",
			args:       args{[]string{"2"}},
			wantW:      "",
			dictionary: nil,
			assertion:  assert.Error,
		},
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
		{
			name:  "usage test name",
			wantW: fmt.Sprintf("%s: print number converted to words\nusage: %s <number>\n", name, name)},
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
