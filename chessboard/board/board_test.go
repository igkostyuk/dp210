package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	type args struct {
		height      int
		width       int
		blackSymbol rune
		whiteSymbol rune
	}
	tests := []struct {
		name      string
		args      args
		want      *Board
		assertion assert.ErrorAssertionFunc
	}{

		{
			"valid args",
			args{height: 1, width: 2, blackSymbol: '*', whiteSymbol: ' '},
			&Board{Height: 1, Width: 2, Squares: [][]rune{{'*', ' '}}}, assert.NoError,
		},
		{
			"negative height",
			args{height: -1, width: 2, blackSymbol: '*', whiteSymbol: ' '},
			nil, assert.Error,
		},
		{
			"negative width",
			args{height: 1, width: -2, blackSymbol: '*', whiteSymbol: ' '},
			nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBoard(tt.args.height, tt.args.width, tt.args.blackSymbol, tt.args.whiteSymbol)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBoard_String(t *testing.T) {
	type fields struct {
		height  int
		width   int
		squares [][]rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"valid field", fields{height: 1, width: 2, squares: [][]rune{{'*', ' '}, {' ', '*'}}}, "* \n *\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &Board{
				Height:  tt.fields.height,
				Width:   tt.fields.width,
				Squares: tt.fields.squares,
			}
			assert.Equal(t, tt.want, br.String())
		})
	}
}

func Test_createSquares(t *testing.T) {
	b, w := '*', ' '

	type args struct {
		height      int
		width       int
		blackSymbol rune
		whiteSymbol rune
	}
	tests := []struct {
		name string
		args args
		want [][]rune
	}{
		{
			"valid args",
			args{height: 6, width: 4, blackSymbol: b, whiteSymbol: w},
			[][]rune{
				{b, w, b, w},
				{w, b, w, b},
				{b, w, b, w},
				{w, b, w, b},
				{b, w, b, w},
				{w, b, w, b},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, createSquares(tt.args.height, tt.args.width, tt.args.blackSymbol, tt.args.whiteSymbol))
		})
	}
}
