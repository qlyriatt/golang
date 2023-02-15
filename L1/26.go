package main

import (
	"fmt"
	"strings"
)

func convertUp(val rune) rune {
	return []rune(strings.ToUpper(string(val)))[0]
}

func convertDown(val rune) rune {
	return []rune(strings.ToLower(string(val)))[0]
}

func checkUnique(s string) bool {
	m := make(map[rune]struct{})

	for _, val := range s {

		// каждый символ - ключ в map
		// если при прохождении по строке подобный ключ (в т.ч. в другом регистре)
		// уже встречался, возвращается false, иначе true
		if _, ok := m[val]; ok {
			return false
		} else if _, ok := m[convertUp(val)]; ok {
			return false
		} else if _, ok := m[convertDown(val)]; ok {
			return false
		}

		m[val] = struct{}{}
	}

	return true
}

func main() {

	s := []string{"aewiophlmntr", "abCdefAaf", "abchdfkoo"}

	for _, val := range s {
		fmt.Printf("%v: %v\n", val, checkUnique(val))
	}
}
