package ecc

import (
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
