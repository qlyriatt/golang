package main

import (
	"fmt"
	"time"
)

func or(signals ...<-chan any) <-chan any {
	ch := make(chan any)
	done := make(chan struct{})

	go func() {
		for _, s := range signals {
			s := s
			go func(done chan struct{}) {
				select {
				case <-done:
					return
				case <-s:
					close(done)
					return
				}
			}(done)
		}
		<-done
		ch <- struct{}{}
	}()
	return ch
}

func signal(after time.Duration) <-chan any {
	ch := make(chan any)
	go func() {
		defer close(ch)
		time.Sleep(after)
	}()
	return ch
}

func main() {

	start := time.Now()

	<-or(
		signal(1*time.Hour),
		signal(4*time.Minute),
		signal(2*time.Second),
		signal(5*time.Second),
	)
	fmt.Printf("elapsed: %v\n", time.Since(start))
}
