package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y int
}

func New(x, y int) Point {
	return Point{x, y}
}

func getDistance(p1 Point, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p1.x-p2.x), 2) + math.Pow(float64(p1.y-p2.y), 2))
}

func main() {
	p1 := New(50, 100)
	p2 := New(100, 70)

	fmt.Printf("%.3f", getDistance(p1, p2))
}
