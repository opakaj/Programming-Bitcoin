package main

import (
	"fmt"

	ecc "github.com/opakaj/ch1/ecc"
)

func main() {
	a := ecc.NewFieldElement(21, 31)
	b := ecc.NewFieldElement(17, 31)

	fmt.Println(a.Add(b))

}
