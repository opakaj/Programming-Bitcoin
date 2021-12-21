package main

import (
	"fmt"

	ecc "github.com/opakaj/ch2/ecc"
)

func main() {
	p1 := ecc.NewPoint(-1, -1, 5, 7)
	p2 := ecc.NewPoint(-1, 1, 5, 7)
	//inf := NewPoint(0,0, 5, 7)

	fmt.Println(p1.Ne(p2))

}
