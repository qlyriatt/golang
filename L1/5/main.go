package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 3 // время работы в секундах

func write(ch chan int, done chan bool, quit chan bool) {
	for {
		if len(done) != 0 { // ждем сигнал о превышении времени работы от checkTime
			close(ch)
			quit <- true // отправляем сигнал завершения в main
			return
		}
		time.Sleep(time.Millisecond * 200)
		rand.Seed(time.Now().UnixMilli())
		ch <- rand.Intn(100)
	}
}

func read(ch chan int) {
	for {
		val, ok := <-ch
		if !ok { // при закрытии канала завершаем функцию
			return
		}
		fmt.Println(val)
	}
}

func checkTime(startTime int64, done chan bool) {
	for {
		if time.Now().Unix()-startTime > N {
			done <- true // при превышении времени работы отправляем сигнал в write
			return
		}
	}
}

func main() {
	startTime := time.Now().Unix()

	ch := make(chan int, 1)
	done := make(chan bool, 1)
	quit := make(chan bool, 1)

	go checkTime(startTime, done)
	go write(ch, done, quit)
	go read(ch)

	<-quit
}
