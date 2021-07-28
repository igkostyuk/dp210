package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/igkostyuk/dp210/envelopes/envelope"
	"github.com/stretchr/testify/assert"
)

func Test_getSizePairValues(t *testing.T) {
	type args struct {
		r  *bufio.Reader
		sp *sizePair
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid values",
			args{
				bufio.NewReader(strings.NewReader("1\n2\n")),
				&sizePair{{name: "TE"}, {name: "ST"}},
			}, "TE: ST: ", assert.NoError,
		},
		{
			"invalid reader",
			args{
				bufio.NewReader(iotest.ErrReader(errors.New("test"))),
				&sizePair{{name: "A"}, {name: "ST"}},
			}, "A: ", assert.Error,
		},
		{
			"invalid second value",
			args{
				bufio.NewReader(strings.NewReader("1\nINVALID\n")),
				&sizePair{{name: "C"}, {name: "D"}},
			}, "C: D: ", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.assertion(t, getSizePairValues(tt.args.r, w, tt.args.sp))
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_getEnvelops(t *testing.T) {
	type args struct {
		r   *bufio.Reader
		sps []sizePair
	}
	tests := []struct {
		name      string
		args      args
		want      []*envelope.Envelope
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid values",
			args{
				bufio.NewReader(strings.NewReader("1\n2\n3\n4\n")),
				[]sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}},
			}, []*envelope.Envelope{
				{Name: "AB", Height: 1, Width: 2},
				{Name: "CD", Height: 3, Width: 4},
			}, "Enter AB envelope sizes\nA: B: Enter CD envelope sizes\nC: D: ",
			assert.NoError,
		},
		{
			"invalid reader",
			args{
				bufio.NewReader(iotest.ErrReader(errors.New("test"))),
				[]sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}},
			}, nil, "Enter AB envelope sizes\nA: ", assert.Error,
		},
		{
			"negative firs reader value",
			args{
				bufio.NewReader(strings.NewReader("-1\n1\n")),
				[]sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}},
			}, nil, "Enter AB envelope sizes\nA: B: ", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := getEnvelops(tt.args.r, w, tt.args.sps)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_getFitEnvelops(t *testing.T) {
	type args struct {
		r   *bufio.Reader
		sps []sizePair
	}
	tests := []struct {
		name      string
		args      args
		want      []*envelopsPair
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid values",
			args{
				bufio.NewReader(strings.NewReader("5\n6\n4\n5\n")),
				[]sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}},
			}, []*envelopsPair{{
				{Name: "CD", Height: 4, Width: 5},
				{Name: "AB", Height: 5, Width: 6},
			}}, "Enter AB envelope sizes\nA: B: Enter CD envelope sizes\nC: D: ",
			assert.NoError,
		},
		{
			"invalid reader",
			args{
				bufio.NewReader(iotest.ErrReader(errors.New("test"))),
				[]sizePair{{{name: "A"}, {name: "B"}}, {{name: "C"}, {name: "D"}}},
			}, nil, "Enter AB envelope sizes\nA: ",
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := getFitEnvelops(tt.args.r, w, tt.args.sps)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func Test_confirm(t *testing.T) {
	type args struct {
		r        *bufio.Reader
		confirms []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"confirm",
			args{bufio.NewReader(strings.NewReader("test\n")), []string{"test"}}, true,
		},
		{
			"not confirm",
			args{bufio.NewReader(strings.NewReader("test\n")), []string{"no"}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, confirm(tt.args.r, tt.args.confirms))
		})
	}
}

func TestTask(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			"AB fit",
			args{
				strings.NewReader("1\n" + "1\n" + "2\n" + "2\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"envelope AB(1.00,1.00) can fit in CD(2.00,2.00)\n" +
				"continue [y yes] ?:",
		},
		{
			"CD fit",
			args{
				strings.NewReader("2\n" + "2\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"envelope CD(1.00,1.00) can fit in AB(2.00,2.00)\n" +
				"continue [y yes] ?:",
		},
		{
			"cant fit",
			args{
				strings.NewReader("1\n" + "1\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"Enter CD envelope sizes\n" +
				"C: " +
				"D: " +
				"envelops can't fit\n" +
				"continue [y yes] ?:",
		},
		{
			"invalid first envelope param",
			args{
				strings.NewReader("INVALID\n" + "1\n" + "1\n" + "1\n" + "no\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"get fit envelopes:parsing size:strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:",
		},
		{
			"invalid second param and continue",
			args{
				strings.NewReader("1\n" + "INVALID\n" + "y\n" + "1\n" + "INVALID\n"),
			}, "Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"get fit envelopes:parsing size:strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:" +
				"Enter AB envelope sizes\n" +
				"A: " +
				"B: " +
				"get fit envelopes:parsing size:strconv.ParseFloat: parsing \"INVALID\": invalid syntax\n" +
				"continue [y yes] ?:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			Task(tt.args.r, w)
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
		{"usage", "test: checks if one envelope can fit in another\n"},
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
		{"main run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
