package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	const N = 20

	var set1 []int
	var set2 []int
	var intersection []int

	// заполнение множеств неодинаковыми числами
	rand.Seed(time.Now().Unix())
	for i := 0; i < N; i++ {
		tmp := rand.Intn(100)
		for _, val := range set1 {
			if val == tmp {
				i--
				break
			}
		}
		set1 = append(set1, tmp)
	}
	for i := 0; i < N; i++ {
		tmp := rand.Intn(100)
		for _, val := range set2 {
			if val == tmp {
				i--
				break
			}
		}
		set2 = append(set2, tmp)
	}

	// проход по всем значениям одного множества и поиск их в другом
	for _, val1 := range set1 {
		for _, val2 := range set2 {
			if val1 == val2 {
				intersection = append(intersection, val1)
				break
			}
		}
	}

	fmt.Println(set1)
	fmt.Println(set2)
	fmt.Println(intersection)
}
