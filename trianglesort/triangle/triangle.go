package triangle

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrSizeSyntax indicates that a value does not have the right syntax for the size type.
	ErrSizeSyntax = errors.New("size should be positive float")
	// ErrSizeRatio indicates that triangle size ratio is wrong.
	ErrSizeRatio = errors.New("not valid sizes: expected a + b <= c or a + c <= b or b + c <= a")
)

// Triangle represent triangle.
type Triangle struct {
	Name    string
	A, B, C float64
}

// NewTriangle create triangle with name and sizes.
func NewTriangle(name string, a, b, c float64) (*Triangle, error) {
	if a <= 0 || b <= 0 || c <= 0 {
		return nil, ErrSizeSyntax
	}
	if a+b <= c || a+c <= b || b+c <= a {
		return nil, ErrSizeRatio
	}
	return &Triangle{Name: name, A: a, B: b, C: c}, nil
}

// Area return triangle area.
func (t *Triangle) Area() float64 {
	// Heron's formula
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// String represent triangle as string.
func (t *Triangle) String() string {
	return fmt.Sprintf("[Triangle %s]: %.2f cm", t.Name, t.Area())
}
