package main

import "fmt"

func reverse(s string) string {
	c := []rune(s)
	n := len(c)

	for i := 0; i < n/2; i++ {
		c[i], c[n-1-i] = c[n-1-i], c[i]
	}

	return string(c)
}

func simpleReverse(s string) string {
	var r string
	for _, c := range s {
		r = string(c) + r
	}

	return r
}

func main() {
	s := "hello world"

	fmt.Println(reverse(s))
	fmt.Println(reverse(s))
}
