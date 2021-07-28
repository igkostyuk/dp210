package triangle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTriangle(t *testing.T) {
	type args struct {
		name string
		a    float64
		b    float64
		c    float64
	}
	tests := []struct {
		name      string
		args      args
		want      *Triangle
		assertion assert.ErrorAssertionFunc
	}{
		{
			"valid triangle",
			args{name: "test", a: 1, b: 1, c: 1},
			&Triangle{Name: "test", A: 1, B: 1, C: 1}, assert.NoError,
		},
		{
			"negative size",
			args{name: "test", a: -1, b: 1, c: 1}, nil, assert.Error,
		},
		{
			"too big size",
			args{name: "test", a: 1, b: 5, c: 1}, nil, assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTriangle(tt.args.name, tt.args.a, tt.args.b, tt.args.c)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTriangle_Area(t *testing.T) {
	type fields struct {
		name string
		a    float64
		b    float64
		c    float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"same sizes #1", fields{a: 5, b: 4, c: 3}, 6},
		{"same sizes #2", fields{a: 4, b: 5, c: 3}, 6},
		{"same sizes #3", fields{a: 3, b: 4, c: 5}, 6},
		{"all sizes 1", fields{a: 1, b: 1, c: 1}, 0.4330127018922193},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Triangle{
				Name: tt.fields.name,
				A:    tt.fields.a,
				B:    tt.fields.b,
				C:    tt.fields.c,
			}
			assert.Equal(t, tt.want, tr.Area())
		})
	}
}

func TestTriangle_String(t *testing.T) {
	type fields struct {
		name string
		a    float64
		b    float64
		c    float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"round area", fields{name: "test", a: 3, b: 4, c: 5}, "[Triangle test]: 6.00 cm"},
		{"float area", fields{name: "test", a: 1, b: 1, c: 1}, "[Triangle test]: 0.43 cm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Triangle{
				Name: tt.fields.name,
				A:    tt.fields.a,
				B:    tt.fields.b,
				C:    tt.fields.c,
			}
			assert.Equal(t, tt.want, tr.String())
		})
	}
}
