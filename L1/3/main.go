package main

import "fmt"

var nums = []int{2, 4, 6, 8, 10}

// отправляем квадрат числа в канал
func getSquare(n int, ch chan int) {
	sq := n * n
	ch <- sq
}

func main() {

	ch := make(chan int, len(nums))

	for _, n := range nums {
		go getSquare(n, ch)
	}

	// суммируем все квадраты, накопившиеся в канале
	sum := 0
	for i := 0; i < cap(nums); i++ {
		sum += <-ch
	}

	fmt.Printf("SqSum: %v", sum)
}
