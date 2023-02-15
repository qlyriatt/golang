package main

import (
	"errors"
	"fmt"
)

func delete(s []int, i int) ([]int, error) {
	if !(i > 0 && i < len(s)) {
		return nil, errors.New("Invalid index")
	}
	tmp := make([]int, 0)
	tmp = append(tmp, s[:i]...)   // забираем часть массива до элемента с индексом i
	tmp = append(tmp, s[i+1:]...) // после него
	return tmp, nil
}

func main() {

	var s = []int{0, 1, 2, 3, 4, 5, 6}

	l := len(s) - 1
	for l > 0 {

		var err error
		s, err = delete(s, l)
		// ??
		// s, err := delete(s, l)
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Println(s)
		}
		l--
	}

}

// !!!
