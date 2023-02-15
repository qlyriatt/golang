package main

import (
	"fmt"
	"math/big"
)

const a, b = 2 << 21, (2 << 30) - 1

func main() {
	var x, y, result big.Int

	x.SetInt64(a)
	y.SetInt64(b)

	result.Add(&x, &y)
	fmt.Println("Add: ", result.String())
	result.SetInt64(0)

	result.Sub(&x, &y)
	fmt.Println("Sub: ", result.String())
	result.SetInt64(0)

	result.Mul(&x, &y)
	fmt.Println("Mul: ", result.String())
	result.SetInt64(0)

	result.Div(&x, &y)
	fmt.Println("Div: ", result.String())
}
