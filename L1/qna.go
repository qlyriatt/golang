package main

import (
	"fmt"
	"strings"
)

// //// 2.
type cat struct{}
type dog struct{}

func (cat) speak() string {
	return "meow"
}

func (dog) speak() string {
	return "woof"
}

type speaker interface {
	speak() string
}

func returnAny(i interface{}) interface{} {
	return i
}

//////

func main() {
	// 1. strings.Builder
	var s1, s2 = "hello ", "world"
	var b strings.Builder
	b.WriteString(s1)
	b.WriteString(s2)
	fmt.Println(b.String(), "\n")

	// 2. Полиморфизм, передача любых параметров в функцию
	var sp speaker = cat{}
	fmt.Println(sp.speak())
	sp = dog{}
	fmt.Println(sp.speak())
	fmt.Println(returnAny(40))
	fmt.Println(returnAny("string"), "\n")

	// 3. RWMutex поддерживает несколько читателей одновременно

	// 4. Запись в буф. канал блокирует только при его заполненности, чтение - при отсутствии значений
	//	  Запись/чтение в небуф. канал блокируют всегда

	// 5. 0

	// 6. Нет

	// 7. ключ: значение

	// 8. make возвращает map, slice, channel, new возвращает указатель на объект не встроенного типа,
	// память заполняется 0

	// 9.
	var m map[string]int
	m = *new(map[string]int)
	m = make(map[string]int)
	m = map[string]int{"key1": 1}
	fmt.Println(m)

	var sl []int
	sl = *new([]int) // не создает массив
	sl = make([]int, 0)
	sl = []int{0, 1, 2, 3, 4}
	sl = sl[1:3]
	fmt.Println(sl, "\n")

	// 10. В функцию передается копия указателя

	// 11. Программа выведет от 0 до 4 в случайном порядке и остановится на Wait(),
	// поскольку счетчик оригинальной WaitGroup не уменьшится из-за передачи по значению в функцию

	// 12. Выведет 0, две переменные находятся в разных блоках видимости

	// 13. [100, 2, 3, 4, 5], обращения в слайс через [] изменяют оригинальный массив,
	// append в данном случае создает копию оригинального массива

	// 14. [b, b, a] [a, a], append создает копию массива

}
