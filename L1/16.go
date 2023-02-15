package main

import (
	"fmt"
	"math/rand"
	"time"
)

func quickSort(a []int) {
	if len(a) < 2 {
		return
	}

	begin := 0
	end := len(a) - 1

	pi := sortPart(a, begin, end)

	quickSortN(a, begin, pi-1)
	quickSortN(a, pi, end)
}

// каждую итерацию получается массив, где все элементы <= опорного расположены слева, остальные справа
func quickSortN(a []int, begin int, end int) {
	if end-begin < 2 {
		return
	}

	pi := sortPart(a, begin, end)
	quickSortN(a, begin, pi-1) // сортируем левую часть
	quickSortN(a, pi, end)     // правую часть
}

func sortPart(a []int, begin int, end int) int {
	p := a[end]
	pi := end

	for i := end - 1; i >= begin; i-- {
		if a[i] > p {
			tmp := a[i]
			for j := i; j < pi; j++ {
				a[j] = a[j+1]
			}
			a[pi] = tmp
			pi--
		}
	}

	return pi
}

const max = 20

func main() {
	var values []int

	rand.Seed(time.Now().Unix())
	i := 0
	for i < max {
		values = append(values, rand.Intn(100))
		i++
	}

	quickSort(values)
	fmt.Print(values)
}
