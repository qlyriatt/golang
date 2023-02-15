package main

import "strings"

const s = "snow dog sun"

func main() {
	words := strings.Fields(s) // получение отдельных слов из строки

	var reverse string
	for i := 0; i < len(words)/2; i++ { // разворот массива слов
		j := len(words) - i - 1
		words[i], words[j] = words[j], words[i]
	}

	for _, word := range words { // объединение слов в строку
		reverse += word + " "
	}

	print(reverse)
}
