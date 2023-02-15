package main

import (
	"errors"
	"fmt"
)

func setBit(num int64, bit int, val int) (int64, error) {

	if bit > 63 || bit < 0 {
		return 0, errors.New("Wrong bit")
	}
	if !(val == 1 || val == 0) {
		return 0, errors.New("Wrong bit value")
	}

	//             bitMask
	// 00000000000000 1 000000000000000
	//            highMask
	// 11111111111111 1 000000000000000
	//             lowMask
	// 00000000000000 0 111111111111111
	// -----high-----   ------low------
	// 01000110110011 0 100110011100110
	//				  |
	// 				 bit

	var empty int64 = 0
	var lowMask int64 = (1 << bit) - 1
	var highMask int64 = ^lowMask
	var bitMask int64 = 1 << bit
	low := num & lowMask   // копия нижних битов num
	high := num & highMask // копия верхних битов num

	if val == 0 {
		empty = num >> (bit + 1)
		empty = empty << (bit + 1)
		empty = empty | low
	} else {
		empty = (high ^ bitMask) | low
	}

	return empty, nil
}

var N int64 = 9223372036854775807

func main() {
	if num, err := setBit(N, 62, 0); err != nil {
		fmt.Print(err)
	} else {
		fmt.Printf("N: %v, new N: %v\n", N, num)
	}
}
