package main

import "fmt"

type Human struct {
	i int
	s string
	f float32
	b bool
}

func (h Human) getInt() int {
	return h.i
}

func (h Human) getString() string {
	return h.s
}

func (h Human) getFloat() float32 {
	return h.f
}

func (h Human) getBool() bool {
	return h.b
}

type Action struct {
	Human
}

func main() {
	a := Action{
		Human{25, "hello", 36.6, true},
	}

	// использование методов Human через Action
	i := a.getInt()
	s := a.getString()
	f := a.getFloat()
	b := a.getBool()

	fmt.Printf("int: %v, string: %v, float: %v, bool: %v", i, s, f, b)
}
