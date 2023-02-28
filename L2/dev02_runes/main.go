package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func Unpack(input string) (string, error) {

	// пустая строка
	if len(input) < 1 {
		return input, nil
	}

	if unicode.IsDigit([]rune(input)[0]) {
		return "", errors.New("первый символ строки - цифра")
	}

	output := []rune{}
	esc := false  // флаг начала escape последовательности
	mult := false // флаг нахождения цифры последним считанным символом
	for _, char := range input {

		if unicode.IsDigit(char) { // цифры

			if mult {
				return "", errors.New("несколько цифр подряд")
			}

			// escape для цифры
			if esc {
				esc = false
			} else {
				// повторить последний символ в выводе num - 1 раз
				last := output[len(output)-1]
				num, err := strconv.Atoi(string(char))
				if err != nil {
					return "", err
				}
				for num > 1 {
					output = append(output, last)
					num--
				}
				mult = true
				continue
			}

		} else if char == '\\' {
			// начало escape
			if !esc {
				esc = true
				mult = false
				continue
			}
			esc = false

		} else {
			if esc {
				return "", errors.New("встречена неверная escape-последовательность")
			}
		}

		// введенный символ добавляется в вывод, кроме случаев escape и повторения
		output = append(output, char)
		mult = false
	}

	if esc {
		return "", errors.New("обнаружен лишний '\\'")
	}
	return string(output), nil
}

func main() {
	var input string
	fmt.Scanln(&input)

	if output, err := Unpack(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(output)
	}
}
