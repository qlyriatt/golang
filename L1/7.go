package main

import (
	"fmt"
	"sync"
)

// удаляем сигнал из lock, блокируя запись для остальных рутин
// по окончании записи возвращаем сигнал
func write(m map[string]int, lock chan struct{}, s string, n int) {
	select {
	case <-lock: // если в lock нет сигнала, то пропускаем запись
		m[s] = n
		lock <- struct{}{}
	default:
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	m := make(map[string]int)
	// если в канале есть сигнал, то m готова для записи
	lock := make(chan struct{}, 1)
	lock <- struct{}{} // отправляем сигнал о доступности m

	// конкурентная запись
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write(m, lock, "i"+fmt.Sprint(i), i)
	}

	wg.Wait()                // ждем завершения всех рутин
	fmt.Printf("m: %v\n", m) // вывод m
}
