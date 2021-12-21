package ecc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
)

var A int = 0
var B int = 7
var P int = int(math.Pow(float64(2), float64(256)) - math.Pow(float64(2), float64(32)) - 977)
var N int64 = 0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141

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

func (p *Point) Add(other *Point) interface{} {
	if !reflect.DeepEqual(p.a, other.a) || !reflect.DeepEqual(p.b, other.b) {
		panic(errors.New("exception, points are not on the curve"))
	}
	//Case 0.0: self is the point at infinity, return other
	if p.x == 0 {
		return other
	}
	//Case 0.1: other is the point at infinity, return self
	if other.x == 0 {
		return p
	}
	//Case 1: self.x == other.x, self.y != other.y
	//Result is point at infinity
	if reflect.DeepEqual(p.x, other.x) && !reflect.DeepEqual(p.y, other.y) {
		return NewPoint(0, 0, p.a, p.b)
	}
	//Case 2: self.x ≠ other.x
	//Formula (x3,y3)==(x1,y1)+(x2,y2)
	//s=(y2-y1)/(x2-x1)
	//x3=s**2-x1-x2
	//y3=s*(x1-x3)-y1
	if !reflect.DeepEqual(p.x, other.x) {
		s := (other.y - p.y) / (other.x - p.x)
		x := math.Pow(float64(s), 2) - float64(p.x) - float64(other.x)
		y := float64(s)*(float64(p.x)-x) - float64(p.y)
		return NewPoint(int(x), int(y), p.a, p.b)
	}
	//Case 4: if we are tangent to the vertical line,
	//we return the point at infinity
	//note instead of figuring out what 0 is for each type
	//we just use 0 * self.x
	if p.Eq(other) && reflect.DeepEqual(p.y, 0*int(p.x)) {
		return NewPoint(0, 0, p.a, p.b)
	}
	//Case 3: self == other
	//Formula (x3,y3)=(x1,y1)+(x1,y1)
	//s=(3*x1**2+a)/(2*y1)
	//x3=s**2-2*x1
	//y3=s*(x1-x3)-y1
	if p.Eq(other) {
		s := float64(3*math.Pow(float64(p.x), float64(2))+float64(p.a)) / float64(2*int(p.y))
		x := math.Pow(s, float64(2)) - float64(2*int(p.x))
		y := s*(float64(p.x)-x) - float64(p.y)
		return NewPoint(int(x), int(y), p.a, p.b)
	}
	return NewPoint(0, 0, 0, 0)
}

func (p *Point) Rmul(coefficient int) *Point {
	coef := mod(int64(coefficient), N)
	//coef := coefficient
	current := p
	result := NewPoint(0, 0, p.a, p.b)
	for coef != 0 {
		if int(coef)&1 != 0 {
			result.Add(current)
		}
		current.Add(current)
		coef >>= 1
	}
	return result
}

type S256Field struct {
	FieldElement
}

//removed prime as an argument
func NewS256Field(num int) (s *S256Field) {
	s = new(S256Field)
	s.num = num
	prime := P
	s.prime = prime
	NewFieldElement(num, prime)
	return
}

func (f *S256Field) Repr() string {
	str := fmt.Sprintf("%064d", (f.num))
	return str
}

/*
func (f *S256Field) Sqrt() int {
	return int(math.Pow(float64(f), float64(((P + 1) / 4))))
}
*/
type S256Point struct {
	Point
	a  int
	b  int
	xx *S256Field //we need to access num variable so we use new xx variable instead of the x variable from Point
	y  int
	//S256Field
}

func NewS256Point(xx, y, a, b interface{}) (self *S256Point) {
	//var a interface{}
	//var b interface{}
	self = new(S256Point)
	a, b = NewS256Field(A), NewS256Field(B)
	if reflect.TypeOf(xx).Kind() == reflect.Int && xx != nil {
		NewS256Point(xx, NewS256Field(y.(int)), a.(int), b.(int))
	} else {
		NewS256Point(xx, y, a, b)
	}
	return
}

func (s256 *S256Point) Rmul2(coefficient int) *S256Point {
	coef := mod(int64(coefficient), N)
	//coef := coefficient
	current := s256
	result := NewS256Point(0, 0, s256.a, s256.b)
	for coef != 0 {
		if int(coef)&1 != 0 {
			result.SAdd(current)
		}
		current.SAdd(current)
		coef >>= 1
	}
	return result
}

func (p *Point) SEq(other *S256Point) bool {
	return reflect.DeepEqual(p.x, other.x) && reflect.DeepEqual(p.y, other.y) && reflect.DeepEqual(p.a, other.a) && reflect.DeepEqual(p.b, other.b)
}

func (p *S256Point) SAdd(other *S256Point) *S256Point {
	if !reflect.DeepEqual(p.a, other.a) || !reflect.DeepEqual(p.b, other.b) {
		panic(errors.New("exception, points are not on the curve"))
	}
	//Case 0.0: p is the point at infinity, return other
	if p.xx == nil {
		return other
	}
	//Case 0.1: other is the point at infinity, return p
	if other.xx == nil {
		return p
	}
	//Case 1: p.x == other.x, p.y != other.y
	//Result is point at infinity
	if reflect.DeepEqual(p.x, other.x) && !reflect.DeepEqual(p.y, other.y) {
		return NewS256Point(0, 0, p.a, p.b)
	}
	//Case 2: p.x ≠ other.x
	//Formula (x3,y3)==(x1,y1)+(x2,y2)
	//s=(y2-y1)/(x2-x1)
	//x3=s**2-x1-x2
	//y3=s*(x1-x3)-y1
	if !reflect.DeepEqual(p.x, other.x) {
		s := (other.y - p.y) / (other.xx.num - p.x)
		x := math.Pow(float64(s), 2) - float64(p.x) - float64(other.x)
		y := float64(s)*(float64(p.x)-x) - float64(p.y)
		return NewS256Point(NewS256Field(int(x)), NewS256Field(int(y)), p.a, p.b)
	}
	//Case 4: if we are tangent to the vertical line,
	//we return the point at infinity
	//note instead of figuring out what 0 is for each type
	//we just use 0 * p.x
	if p.SEq(other) && reflect.DeepEqual(p.y, 0*int(p.x)) {
		return NewS256Point(0, 0, p.a, p.b)
	}
	//Case 3: p == other
	//Formula (x3,y3)=(x1,y1)+(x1,y1)
	//s=(3*x1**2+a)/(2*y1)
	//x3=s**2-2*x1
	//y3=s*(x1-x3)-y1
	if p.SEq(other) {
		s := float64(3*math.Pow(float64(p.x), float64(2))+float64(p.a)) / float64(2*int(p.y))
		x := math.Pow(s, float64(2)) - float64(2*int(p.x))
		y := s*(float64(p.x)-x) - float64(p.y)
		return NewS256Point(NewS256Field(int(x)), NewS256Field(int(y)), p.a, p.b)
	} else {
		return NewS256Point(0, 0, 0, 0)
	}

}

func (s256 *S256Point) verify(z int64, sig *Signature) bool {
	s_inv := mod(int64(math.Pow(sig.s.(float64), float64(N-2))), N)
	u := mod((z * s_inv), N)
	v := mod((sig.r.(int64) * s_inv), N)
	total := new(S256Point)
	total = (G.Rmul2(int(u))).SAdd(s256.Rmul2(int(v))) //nolonger using mul fxn
	return total.xx.num == sig.r
}

type Signature struct {
	r interface{}
	s interface{}
}

func NewSignature(r interface{}, s interface{}) (ss *Signature) {
	ss = new(Signature)
	ss.r = r
	ss.s = s
	return
}

func (S *Signature) Repr() string {
	str := fmt.Sprintf("Signature(%s,%s)", S.r, S.s)
	return str
}

var G = NewS256Point(0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798, 0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8, 0, 0)

type PrivateKey struct {
	secret int
	point  interface{}
}

func NewPrivateKey(secret int) (pk *PrivateKey) {
	pk = new(PrivateKey)
	pk.secret = secret
	pk.point = new(S256Point)
	pk.point = (G.Rmul2(secret))
	return
}

/*
func Mul(other int, f *S256Point) *S256Point {
	if f.a != 0 && f.b != 0 {
		panic(fmt.Errorf("TypeError: %v", "Cannot multiply by the G because of non zero a and b"))
	}
	num := NewS256Point(NewS256Field((other * f.x)), NewS256Field((other * f.y)), (other * f.a), (other * f.b))
	return num
}
*/
func (p *PrivateKey) hex() {
	fmt.Printf("%64d", (p.secret))
}

func (pk *PrivateKey) sign(z int64) *Signature {
	k := pk.deterministic_k(z)
	r := (G.Rmul2(int(k))).xx.num
	k_inv := mod(int64(math.Pow(float64(k), float64(N-2))), N)
	s := (z + int64(r*pk.secret)) * k_inv % N
	if float64(s) > float64(N)/float64(2) {
		s = N - s
	}
	return NewSignature(r, s)
}

func (pk *PrivateKey) deterministic_k(z int64) int64 {
	k := func(repeated []byte, n int) (result []byte) {
		for i := 0; i < n; i++ {
			result = append(result, repeated...)
		}
		return result
	}([]byte("\u0000"), 32)
	v := func(repeated []byte, n int) (result []byte) {
		for i := 0; i < n; i++ {
			result = append(result, repeated...)
		}
		return result
	}([]byte("\u0001"), 32)
	if z > N {
		z -= N
	}

	z_bytes := make([]byte, 32)
	binary.BigEndian.PutUint64(z_bytes, uint64(z))
	secret_bytes := make([]byte, 32)
	binary.BigEndian.PutUint64(secret_bytes, uint64(pk.secret))

	kk := hmac.New(sha256.New, k)                                                         // Create a new HMAC by defining the hash type and the key (as byte array)
	kk.Write(append(append(append(v, []byte("\u0000")...), secret_bytes...), z_bytes...)) // Write Data to it
	//sha := hex.EncodeToString(kk.Sum(nil))
	k, _ = hex.DecodeString(hex.EncodeToString(kk.Sum(nil))) //DecodeString returns the bytes represented by the hexadecimal string s.

	vv := hmac.New(sha256.New, k)
	vv.Write(v)
	v, _ = hex.DecodeString(hex.EncodeToString(vv.Sum(nil)))

	kk = hmac.New(sha256.New, v)
	kk.Write(append(append(append(v, []byte("\u0000")...), secret_bytes...), z_bytes...))
	k, _ = hex.DecodeString(hex.EncodeToString(kk.Sum(nil)))

	vv = hmac.New(sha256.New, k)
	vv.Write(v)
	v, _ = hex.DecodeString(hex.EncodeToString(vv.Sum(nil)))

	for {
		vv = hmac.New(sha256.New, k)
		vv.Write(v)
		v, _ = hex.DecodeString(hex.EncodeToString(vv.Sum(nil)))
		candidate := int64(binary.LittleEndian.Uint64(v))
		if candidate >= 1 && candidate < N {
			return candidate
		}
		kk = hmac.New(sha256.New, k)
		kk.Write(append(k, []byte("\u0000")...))
		k, _ = hex.DecodeString(hex.EncodeToString(kk.Sum(nil)))

		vv = hmac.New(sha256.New, k)
		vv.Write(v)
		v, _ = hex.DecodeString(hex.EncodeToString(vv.Sum(nil)))

	}
}
