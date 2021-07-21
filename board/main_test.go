package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestNewBoard(t *testing.T) {
	type args struct {
		height int
		width  int
		white  rune
		black  rune
	}
	tests := []struct {
		name    string
		args    args
		want    *Board
		wantErr bool
	}{
		{"valid args", args{height: 1, width: 2, black: '*', white: ' '}, &Board{height: 1, width: 2, squares: [][]rune{{'*', ' '}}}, false},
		{"negative height", args{height: -1, width: 2, black: '*', white: ' '}, nil, true},
		{"negative width", args{height: 1, width: -2, black: '*', white: ' '}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBoard(tt.args.height, tt.args.width, tt.args.black, tt.args.white)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBoard() = %v, want %v", got, tt.want)
			}

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
				height:  tt.fields.height,
				width:   tt.fields.width,
				squares: tt.fields.squares,
			}
			if got := br.String(); got != tt.want {
				t.Errorf("Board.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createSquares(t *testing.T) {
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
			args{height: 6, width: 4, blackSymbol: '*', whiteSymbol: ' '},
			[][]rune{
				{'*', ' ', '*', ' '},
				{' ', '*', ' ', '*'},
				{'*', ' ', '*', ' '},
				{' ', '*', ' ', '*'},
				{'*', ' ', '*', ' '},
				{' ', '*', ' ', '*'},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSquares(tt.args.height, tt.args.width, tt.args.blackSymbol, tt.args.whiteSymbol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createSquares() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createBoard(t *testing.T) {
	type args struct {
		parameters  []string
		blackSymbol rune
		whiteSymbol rune
	}
	tests := []struct {
		name    string
		args    args
		want    *Board
		wantErr bool
	}{
		{
			"valid args",
			args{[]string{"2", "3"}, '^', ' '},
			&Board{2, 3, [][]rune{{'^', ' ', '^'}, {' ', '^', ' '}}}, false,
		},
		{
			"invalid params len",
			args{[]string{"2", "3", "4"}, '^', ' '}, nil, true,
		},
		{
			"invalid first params",
			args{[]string{"invalid", "4"}, '^', ' '}, nil, true,
		},
		{
			"invalid second params",
			args{[]string{"2", "invalid"}, '^', ' '}, nil, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createBoard(tt.args.parameters, tt.args.blackSymbol, tt.args.whiteSymbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("createBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createBoard() = %v, want %v", got, tt.want)
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
		{"valid args", args{[]string{"2", "8"}},
			"* * * * \n * * * *\n",
			false},
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
		{"usage", "test: print chessboard\nusage: test <height> <width>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			usage(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Task() =\n%v\n, want \n%v\n", gotW, tt.wantW)
			}

		})
	}
}
