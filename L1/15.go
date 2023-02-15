package main

import (
	"fmt"
	"strings"
)

// var justString string
// func someFunc() {
//   v := createHugeString(1 << 10)
//   justString = v[:100]
// }
//
// после завершения функции v останется в памяти до момента освобождения justString

func createHugeString(n int) string {
	var s string
	for n > 0 {
		s += "a"
		n--
	}
	return s // возвращаем созданную строку
}

func getSubstring(sub string) {
	s := createHugeString(1000)
	sub = strings.Clone(s[:100])
	// после выполнения s будет высвобождена из памяти
}

func main() {
	var sub string
	getSubstring(sub)
	fmt.Print(sub)
}
