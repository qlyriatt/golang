package main

import (
	"fmt"
	"sync"
)

func readToChannel(nums chan int) {
	for _, v := range A { // чтение значений из массива в канал
		nums <- v
	}
	close(nums)
	wg.Done()
}

func writeToOut(nums chan int, doubleNums chan int) {
	for val := range nums { // вывод умноженных значений во второй канал
		doubleNums <- val * 2
	}
	close(doubleNums)
	wg.Done()
}

var A = []int{1, 2, 3, 4, 5, 6, 7, 8}
var wg sync.WaitGroup

func main() {
	nums := make(chan int)
	doubleNums := make(chan int)

	wg.Add(2)
	go readToChannel(nums)
	go writeToOut(nums, doubleNums)

	for val := range doubleNums {
		fmt.Println(val)
	}
	wg.Wait()
}
