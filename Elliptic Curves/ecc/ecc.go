package ecc

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
)

type FieldElement struct {
	num   int
	prime int
}

func mod(x, y int64) int64 {
	bx, by := big.NewInt(x), big.NewInt(y)
	return new(big.Int).Mod(bx, by).Int64()
}

//named return value *fieldelement
func NewFieldElement(num int, prime int) (f *FieldElement) {
	f = new(FieldElement)
	if num >= prime || num < 0 {
		error := fmt.Sprintf("Num %d not in field range 0 to %d", num, (prime - 1))
		panic(error)
	}
	f.num = num
	f.prime = prime
	return
}

func (f *FieldElement) Repr() string {
	s := fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.num)
	return s
}

func (f *FieldElement) Eq(other *FieldElement) bool {
	if other == nil {
		return false
	}
	return reflect.DeepEqual(f.num, other.num) && reflect.DeepEqual(f.prime, other.prime)
}

func (f *FieldElement) Ne(other *FieldElement) bool {
	return !reflect.DeepEqual(f.num, other.num) && reflect.DeepEqual(f.prime, other.prime)
}

func (f *FieldElement) Add(other *FieldElement) *FieldElement {
	if !reflect.DeepEqual(f.prime, other.prime) {
		panic(fmt.Errorf("TypeError: %v", "Cannot Add two numbers in different Fields"))
	}
	num := (f.num + other.num) % f.prime
	return NewFieldElement(num, f.prime)
}
func (f *FieldElement) Sub(other *FieldElement) *FieldElement {
	if !reflect.DeepEqual(f.prime, other.prime) {
		panic(fmt.Errorf("TypeError: %v", "Cannot Subtract two numbers in different Fields"))
	}
	num := (f.num - other.num) % f.prime
	return NewFieldElement(num, f.prime)
}

func (f *FieldElement) Mul(other *FieldElement) *FieldElement {
	if !reflect.DeepEqual(f.prime, other.prime) {
		panic(fmt.Errorf("TypeError: %v", "Cannot multiply two numbers in different Fields"))
	}
	num := (f.num * other.num) % f.prime
	return NewFieldElement(num, f.prime)
}

func (f *FieldElement) Pow(exponent int) *FieldElement {
	var n int64 = mod(int64(exponent), int64(f.prime-1))
	//fmt.Println("n is ", n)
	num := int(math.Pow(float64(f.num), float64(n))) % (f.prime)
	//fmt.Println("num is ", num)
	return NewFieldElement(num, f.prime)
}

func (f *FieldElement) Truediv(other *FieldElement) {
	if !reflect.DeepEqual(f.prime, other.prime) {
		panic(fmt.Errorf("TypeError: %v", "Cannot divide two numbers in different Fields"))
	}
	num := f.num * (int(math.Pow(float64(other.num), float64(f.prime-2))) % f.prime) % f.prime
	NewFieldElement(num, f.prime)
}

type Point struct {
	a int
	b int
	x int
	y int
}

func NewPoint(x, y, a, b int) (p *Point) {
	p = new(Point)
	p.a = a
	p.b = b
	p.x = x
	p.y = y

	if p.x == 0 && p.y == 0 {
		return
	}
	if !reflect.DeepEqual(math.Pow(float64(p.y), float64(2)), math.Pow(float64(p.x), float64(3))+float64(a*x)+float64(b)) {
		error := fmt.Sprintf("(%d, %d) is not on the curve", x, y)
		panic(error)
	}
	return
}

func (p *Point) Eq(other *Point) bool {
	return reflect.DeepEqual(p.x, other.x) && reflect.DeepEqual(p.y, other.y) && reflect.DeepEqual(p.a, other.a) && reflect.DeepEqual(p.b, other.b)
}

func (p *Point) Ne(other *Point) bool {
	return !p.Eq(other)
}

func (p *Point) Repr() string {
	if p.x == 0 {
		return "Point(infinity)"
	} else {
		err := fmt.Sprintf("Point(%d,%d)_%d_%d", p.x, p.y, p.a, p.b)
		return err
	}
}

func (p *Point) Add(other *Point) *Point {
	//p = new(Point)
	if !reflect.DeepEqual(p.a, other.a) || !reflect.DeepEqual(p.b, other.b) {
		panic(errors.New("exception, points are not on the curve"))
	}
	if p.x == 0 {
		return other
	}
	if other.x == 0 {
		return p
	}
	if reflect.DeepEqual(p.x, other.x) && !reflect.DeepEqual(p.y, other.y) {
		return NewPoint(0, 0, p.a, p.b)
	}
	if !reflect.DeepEqual(p.x, other.x) {
		s := (other.y - p.y) / (other.x - p.x)
		x := math.Pow(float64(s), 2) - float64(p.x) - float64(other.x)
		y := float64(s)*(float64(p.x)-x) - float64(p.y)
		return NewPoint(int(x), int(y), p.a, p.b)
	}
	if p.Eq(other) && reflect.DeepEqual(p.y, 0*int(p.x)) {
		return NewPoint(0, 0, p.a, p.b)
	}
	if p.Eq(other) {
		s := float64(3*math.Pow(float64(p.x), float64(2))+float64(p.a)) / float64(2*int(p.y))
		x := math.Pow(s, float64(2)) - float64(2*int(p.x))
		y := s*(float64(p.x)-x) - float64(p.y)
		return NewPoint(int(x), int(y), p.a, p.b)
	}
	return NewPoint(0, 0, 0, 0)
}
