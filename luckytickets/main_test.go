package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func Test_countDigits(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"four digits", args{1_000}, 4},
		{"ten digits", args{1_000_000_000}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, countDigits(tt.args.n))
		})
	}
}

func Test_isMoskowLucky(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"101", args{101}, true},
		{"112", args{112}, false},
		{"1515", args{1515}, true},
		{"1415", args{1415}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isMoskowLucky(tt.args.n))
		})
	}
}

func Test_isPiterLucky(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"101", args{101}, false},
		{"112", args{112}, true},
		{"211", args{211}, true},
		{"121211", args{121211}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isPiterLucky(tt.args.n))
		})
	}
}

func Test_countNumbers(t *testing.T) {
	type args struct {
		min         int
		max         int
		countMethod string
	}
	tests := []struct {
		name      string
		args      args
		want      int
		assertion assert.ErrorAssertionFunc
	}{
		{
			"moscow", args{min: 120123, max: 320320, countMethod: "Moskow"},
			11187, assert.NoError,
		},
		{
			"piter", args{min: 120123, max: 320320, countMethod: "Piter"},
			5790, assert.NoError,
		},
		{
			"invalid", args{min: 120123, max: 320320, countMethod: "invalid"},
			0, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countNumbers(tt.args.min, tt.args.max, tt.args.countMethod)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_readTicketNumber(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name      string
		args      args
		want      int
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid number",
			args{bufio.NewReader(strings.NewReader("000001\n"))},
			1, assert.NoError,
		},
		{
			"too short number",
			args{bufio.NewReader(strings.NewReader("00001\n"))},
			0, assert.Error,
		},
		{
			"too long number",
			args{bufio.NewReader(strings.NewReader("1000001\n"))},
			0, assert.Error,
		},
		{
			"invalid",
			args{bufio.NewReader(strings.NewReader("invalid\n"))},
			0, assert.Error,
		},
		{
			"err reader",
			args{bufio.NewReader(iotest.ErrReader(errors.New("test")))},
			0, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTicketNumber(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_readCountingMethod(t *testing.T) {
	method := "test method"
	tempFileName := createTempFileWithText(method)
	defer func() { // clean up
		err := os.Remove(tempFileName)
		if err != nil {
			log.Fatal(err)
		}
	}()
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"temp file",
			args{bufio.NewReader(strings.NewReader(tempFileName + "\n"))},
			"test method", assert.NoError},
		{
			"reader error",
			args{bufio.NewReader(iotest.ErrReader(errors.New("error")))},
			"", assert.Error,
		},
		{
			"invalid file",
			args{bufio.NewReader(strings.NewReader("invalid\n"))},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCountingMethod(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_readLine(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid buffer",
			args{bufio.NewReader(strings.NewReader("valid\n"))},
			"valid", assert.NoError,
		},
		{
			"invalid buffer",
			args{bufio.NewReader(iotest.ErrReader(errors.New("test")))},
			"", assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLine(tt.args.r)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTask(t *testing.T) {
	tempFileMoskow := createTempFileWithText("Moskow")
	defer func() { // clean up
		err := os.Remove(tempFileMoskow)
		if err != nil {
			log.Fatal(err)
		}
	}()

	tempFileNotExist := createTempFileWithText("Not exist")
	defer func() { // clean up
		err := os.Remove(tempFileNotExist)
		if err != nil {
			log.Fatal(err)
		}
	}()

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name      string
		args      args
		wantW     string
		assertion assert.ErrorAssertionFunc
	}{
		{

			"moskow lucky",
			args{
				bufio.NewReader(strings.NewReader(tempFileMoskow + "\n120123\n320320\n")),
			}, "Enter config filename: " +
				"count method: Moskow\n" +
				"Min: " +
				"Max: " +
				"--Result--\n" +
				"11187", assert.NoError,
		},
		{

			"invalid file",
			args{
				bufio.NewReader(strings.NewReader("INVALID\n" + "120123\n320320\n")),
			}, "Enter config filename: ", assert.Error,
		},
		{

			"invalid firs ticket number",
			args{
				bufio.NewReader(strings.NewReader(tempFileMoskow + "\n" + "INVALID\n320320\n")),
			}, "Enter config filename: " +
				"count method: Moskow\n" +
				"Min: ", assert.Error,
		},
		{

			"invalid second number",
			args{
				bufio.NewReader(strings.NewReader(tempFileMoskow + "\n" + "120123\nINVALID\n")),
			}, "Enter config filename: " +
				"count method: Moskow\n" +
				"Min: " +
				"Max: ", assert.Error,
		},
		{

			"invalid method in file",
			args{
				bufio.NewReader(strings.NewReader(tempFileNotExist +
					"\n" + "120123\n" + "320320\n"),
				),
			}, "Enter config filename: " +
				"count method: Not exist\n" +
				"Min: " +
				"Max: ", assert.Error,
		},
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
		{"usage", fmt.Sprintf("%s: counting lucky numbers\n", name)},
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

func createTempFileWithText(text string) string {
	f, err := os.CreateTemp(os.TempDir(), "temp")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(text); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}
