package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// завершается на wg.Done()
func waitGroup(wg *sync.WaitGroup) {
	fmt.Println("WaitGroup")
	wg.Done()
}

// завершается при получении сигнала(значения) по каналу stop
func checkChannel(stop chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Println("checkChannel")
			close(stop)
			return
		default:
			continue
		}
	}
}

// завершается при закрытии канала closed
func closedChannel(closed chan struct{}) {
	for {
		select {
		case <-closed:
			fmt.Println("closedChannel")
			return
		default:
			continue
		}
	}
}

// завершается при закрытии канала для получения данных
func closedReceiverChannel(send chan struct{}) {
	for {
		_, ok := <-send
		if !ok {
			fmt.Println("closedReceiverChannel")
			return
		}
	}
}

// завершается при вызове функции cancel в main
func contextCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("contextCancel")
			return
		default:
			continue
		}
	}
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)
	go waitGroup(&wg)
	wg.Wait() // ожидание завершения всех рутин в wg

	stop := make(chan struct{})
	go checkChannel(stop)
	stop <- struct{}{}

	closed := make(chan struct{})
	go closedChannel(closed)
	close(closed)

	send := make(chan struct{})
	go closedReceiverChannel(send)
	send <- struct{}{}
	close(send)

	ctx, cancel := context.WithCancel(context.Background())
	go contextCancel(ctx)
	cancel()

	time.Sleep(time.Second)
}
