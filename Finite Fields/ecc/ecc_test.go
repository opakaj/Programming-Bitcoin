package ecc

import (
	"testing"
)

/*
func (f *FieldElement) TestNe(t *testing.T) {
	a := NewFieldElement(2, 31)
	b := NewFieldElement(2, 31)
	c := NewFieldElement(15, 31)
	f.assertEqual(a, b)
	f.assertTrue(a.Neq(*c))
	f.assertFalse(a.Neq(*b))
}
*/

func (f *FieldElement) TestAdd(t *testing.T) {

	a := NewFieldElement(2, 31)
	b := NewFieldElement(15, 31)
	c := NewFieldElement(17, 31)

	if c != a.Add(b) {
		t.Errorf("got %q, wanted %q", a.Add(b), c)
	}

	d := NewFieldElement(17, 31)
	e := NewFieldElement(21, 31)
	g := NewFieldElement(7, 31)

	if c != d.Add(e) {
		t.Errorf("got %q, wanted %q", a.Add(b), g)
	}

}

/*
func (f *FieldElement) TestSub(t *testing.T) {
	a := NewFieldElement(29, 3)
	b := NewFieldElement(14, 31)
	assert.Equal(t, a.Sub(b), NewFieldElement(25, 31))
	a = NewFieldElement(15, 31)
	b = NewFieldElement(30, 31)
	assert.Equal(t, a.Sub(b), NewFieldElement(16, 31))
}
*/
