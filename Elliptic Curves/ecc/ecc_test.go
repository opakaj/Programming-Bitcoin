package ecc

import (
	"fmt"
	"testing"
)

func (f *Point) TestNe(t *testing.T) {
	a := NewPoint(3, -7, 5, 7)
	b := NewPoint(18, 77, 5, 7)
	c := true

	if c != a.Ne(b) {
		out := fmt.Sprintf("got %t, wanted %t", a.Ne(b), c)
		t.Fatalf(out)
	}

}

/*func TestAdd0(t *testing.T)
      a = Point(x=None, y=None, a=5, b=7)
      b = Point(x=2, y=5, a=5, b=7)
      c = Point(x=2, y=-5, a=5, b=7)
      self.assertEqual(a + b, b)
      self.assertEqual(b + a, b)
      self.assertEqual(b + c, a)

  def test_add1(self):
      a = Point(x=3, y=7, a=5, b=7)
      b = Point(x=-1, y=-1, a=5, b=7)
      self.assertEqual(a + b, Point(x=2, y=-5, a=5, b=7))

  def test_add2(self):
      a = Point(x=-1, y=-1, a=5, b=7)
      self.assertEqual(a + a, Point(x=18, y=77, a=5, b=7))
*/
