package main

import "fmt"

func xorSwap(a *int, b *int) {
	*a = *a ^ *b
	*b = *b ^ *a
	*a = *a ^ *b
}

func main() {
	var a, b = 10, 5
	xorSwap(&a, &b) // обмен переменных через xor
	fmt.Println(a, " ", b)

	a, b = b, a // обмен присваиванием
	fmt.Println(a, " ", b)
}
