package main

import "fmt"

// вывод числа и его квадрата
func getSquare(n int, ready chan bool) {
	sq := n * n
	fmt.Printf("square of %v is %v\n", n, sq)

	// по готовности отправляем сигнал в канал ready
	ready <- true
}

func main() {

	ready := make(chan bool, 5)
	for n := 2; n <= 10; n += 2 {
		go getSquare(n, ready)
	}

	// после 5 сигналов готовности завершаем main
	for {
		if len(ready) == 5 {
			return
		}
	}
}
