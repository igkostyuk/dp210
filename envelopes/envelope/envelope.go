package envelope

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrSizeSyntax indicates that a value does not have the right syntax for the size type.
	ErrSizeSyntax = errors.New("size should be positive float")
)

// Envelope represent envelope.
type Envelope struct {
	Name   string
	Height float64
	Width  float64
}

// NewEnvelope create envelope with name height width sizes.
func NewEnvelope(name string, height, width float64) (*Envelope, error) {
	if height <= 0 || width <= 0 {
		return nil, ErrSizeSyntax
	}
	return &Envelope{Name: name, Height: height, Width: width}, nil
}

// String return string representation of envelope.
func (e *Envelope) String() string {
	return fmt.Sprintf("%s(%.2f,%.2f)", e.Name, e.Height, e.Width)
}

// IsFitsIn indicate if envelope can fit in argument envelope.
func (e *Envelope) IsFitsIn(fe *Envelope) bool {
	a, b := fe.Width, fe.Height
	q, p := e.Width, e.Height
	// Rectangles in Rectangles John E. Wetzel
	// https://www.researchgate.net/publication/
	// 273573138_Rectangles_in_Rectangles(page 204)
	if b > a {
		a, b = b, a
	}
	if q > p {
		q, p = p, q
	}
	if q > b {
		return false
	}
	if p < a {
		return true
	}
	// fits diagonally
	return b*(p*p+q*q) > (2*p*q*a + (p*p-q*q)*math.Sqrt(p*p+q*q-a*a))
}
