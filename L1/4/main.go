package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func populateChannel(data chan int, done chan bool) {

	for {
		time.Sleep(time.Second * 1)
		if len(done) != 0 {
			close(data)  // закрываем канал
			done <- true // отправляем сигнал о закрытии в main
			return
		}
		if len(data) > 50 {
			continue
		}

		rand.Seed(time.Now().UnixMilli())
		for i := 0; i < 100; i++ {
			if len(done) != 0 {
				return
			}
			data <- rand.Intn(100)
		}
	}

}

func readChannel(data chan int) {
	for {
		time.Sleep(time.Millisecond * 400)
		val, ok := <-data
		if !ok { // выход из функции при закрытом канале
			return
		}
		fmt.Print(val, "\n")
	}
}

const workersNum = 5

func main() {

	data := make(chan int)
	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go populateChannel(data, done)

	for i := 0; i < workersNum; i++ {
		go readChannel(data)
	}

	// ожидание sigint или sigterm
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	done <- true // при перехвате сигнала завершаем populateChannel
	<-done       // завершение программы только после закрытия data
	return
}
