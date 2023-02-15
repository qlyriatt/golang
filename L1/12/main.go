package main

import "fmt"

var strings = []string{"cat", "cat", "dog", "cat", "tree"}

func main() {

	set := make(map[string]struct{})
	for _, string := range strings {
		set[string] = struct{}{}
	}

	// в map в качестве ключей окажутся только уникальные значения - нужное множество
	for key := range set {
		fmt.Print(key, " ")
	}
	fmt.Println()
}
