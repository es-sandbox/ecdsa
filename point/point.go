package point

import (
	"github.com/evgeniy-scherbina/ecdsa/ring_modulo"
)

type EllipticCurveParams struct {
	A              int
	B              int
	Modulo         int
	BasePointOrder int // redundant
	BasePoint      *Point
}

func (params *EllipticCurveParams) Copy() *EllipticCurveParams {
	return &EllipticCurveParams{
		params.A,
		params.B,
		params.Modulo,
		params.BasePointOrder,
		params.BasePoint,
	}
}

var TestEllipticCurveParams = EllipticCurveParams{
	A:              0,
	B:              7,
	Modulo:         67,
	BasePointOrder: 79, // redundant
	BasePoint:      New(2, 22),
}

var modulo = ring_modulo.MakeConstructor(TestEllipticCurveParams.Modulo)

type Point struct {
	x int
	y int
}

func New(x, y int) *Point {
	return &Point{x, y}
}

func origin() *Point {
	return &Point{0, 0}
}

func (p *Point) X() int {
	return p.x
}

/*
func (p *Point) Y() int {
	return p.y
}
*/

// c = (qy — py) / (qx — px)
// rx = c2 — px — qx
// ry = c (px — rx) — py
func (p *Point) Add(q *Point) *Point {
	if p.isEqual(q) {
		return p.AddSelf()
	}

	r := origin()

	c1 := modulo().Set(q.y).Sub(p.y)
	c2 := modulo().Set(q.x).Sub(p.x)
	c := c1.Div(c2.Get()).Get()

	r.x = modulo().Set(c * c).Sub(p.x).Sub(q.x).Get()
	r.y = modulo().Set(p.x - r.x).Mul(c).Sub(p.y).Get()

	return r
}

// c = (3px2 + a) / 2py
// rx = c2 — 2px
// ry = c (px — rx) — py
func (p *Point) AddSelf() *Point {
	r := origin()

	c1 := modulo().Set(3).Mul(p.x).Mul(p.x).Add(TestEllipticCurveParams.A)
	c2 := modulo().Set(2).Mul(p.y)
	c := c1.Div(c2.Get()).Get()

	tmp := modulo().Set(2).Mul(p.x)
	r.x = modulo().Set(c * c).Sub(tmp.Get()).Get()

	r.y = modulo().Set(p.x - r.x).Mul(c).Sub(p.y).Get()

	return r
}

func (p *Point) Copy() *Point {
	return &Point{p.x, p.y}
}

func (p *Point) isEqual(q *Point) bool {
	return p.x == q.x && p.y == q.y
}