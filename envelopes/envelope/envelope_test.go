package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvelope(t *testing.T) {
	type args struct {
		name   string
		height float64
		width  float64
	}
	tests := []struct {
		name      string
		args      args
		want      *Envelope
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid params", args{"AB", 1.1, 2.2},
			&Envelope{Name: "AB", Height: 1.1, Width: 2.2}, assert.NoError,
		},
		{
			"negative first param",
			args{"AB", -1.1, 2.2}, nil, assert.Error,
		},
		{
			"negative second param",
			args{"AB", 1.1, -2.2}, nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnvelope(tt.args.name, tt.args.height, tt.args.width)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnvelope_IsFitsIn(t *testing.T) {
	type fields struct {
		name   string
		height float64
		width  float64
	}
	type args struct {
		br *Envelope
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"size bigger than argument size",
			fields{"ab", 10, 10}, args{&Envelope{"cd", 5, 5}}, false,
		},
		{
			"size smaller than argument size",
			fields{"ab", 5, 5}, args{&Envelope{"cd", 10, 10}}, true,
		},
		{
			"height bigger than argument size",
			fields{"ab", 9, 4}, args{&Envelope{"cd", 5, 10}}, true,
		},
		{
			"width bigger than argument size",
			fields{"ab", 4, 9}, args{&Envelope{"cd", 10, 5}}, true,
		},
		{
			"diagonal fit",
			fields{"ab", 10, 1}, args{&Envelope{"cd", 9, 9}}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Envelope{
				Name:   tt.fields.name,
				Height: tt.fields.height,
				Width:  tt.fields.width,
			}
			assert.Equal(t, tt.want, r.IsFitsIn(tt.args.br))
		})
	}
}

func TestEnvelope_String(t *testing.T) {
	type fields struct {
		name   string
		height float64
		width  float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"int values",
			fields{name: "test", height: 1, width: 1},
			"test(1.00,1.00)",
		},
		{
			"float values",
			fields{name: "test", height: 1.2345, width: 1.2345},
			"test(1.23,1.23)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Envelope{
				Name:   tt.fields.name,
				Height: tt.fields.height,
				Width:  tt.fields.width,
			}
			assert.Equal(t, tt.want, e.String())
		})
	}
}
