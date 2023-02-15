package main

import "sync"

type Counter struct {
	count int
	m     sync.Mutex
}

// блокируем запись мьютексом
func (c *Counter) Add() {
	c.m.Lock()
	c.count++
	c.m.Unlock()
}

func (c *Counter) Get() int {
	return c.count
}

func incrementor(c *Counter) {
	i := 0
	for i < 100 {
		c.Add()
		i++
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	var c Counter
	i := 10

	// 10 рутин по 100 инкрементов
	for i > 0 {
		wg.Add(1)
		go incrementor(&c)
		i--
	}
	wg.Wait()

	print(c.Get()) // всегда == 1000
}
