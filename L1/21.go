package main

import (
	"fmt"
	"math"
)

type Square struct {
	a int
}

func (s Square) getArea() int {
	return s.a * s.a
}

type Circle struct {
	r int
}

func (c Circle) getCirc() float32 {
	return math.Pi * float32(2*c.r)
}

type SquareAdapter struct {
	Square
}

func (sa SquareAdapter) getCirc() int {
	return sa.a * 4
}

func main() {
	s := Square{5}
	c := Circle{5}
	sa := SquareAdapter{s} // адаптер позволяет Square вызывать методы для Circle

	fmt.Println(s.getArea())
	fmt.Println(c.getCirc())
	fmt.Println(sa.getCirc())
}
