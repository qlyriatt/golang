package main

import (
	"fmt"
	"strings"
)

// структура сообщения
type message struct {
	s            string
	isCoded      bool
	timeToDecode int
}

func acceptMessage(m message) string {
	if m.isCoded {
		return decodeMessage(m)
	}
	return "acceptMessage: " + m.s
}

func decodeMessage(m message) string {
	if m.timeToDecode > 3 {
		return decodeHard(m)
	}

	// decode
	m.s = strings.ReplaceAll(m.s, "1", "l")
	return "decodeMessage: " + m.s
}

func decodeHard(m message) string {
	// decode "smart"
	m.s = strings.ReplaceAll(m.s, "3", "e")
	m.s = strings.ReplaceAll(m.s, "1", "ll")
	m.s = strings.ReplaceAll(m.s, "0", "o")
	return "decodeHard: " + m.s
}

func main() {
	uncodedMsg := message{"hello", false, 0}
	codedMsg := message{"he11o", true, 2}
	hardMsg := message{"h310", true, 5}
	var messages = []message{uncodedMsg, codedMsg, hardMsg}

	for _, m := range messages {
		// сообщение будет обработано в одной из функций
		//
		// при невозможности обработки функции передают сообщения
		// в следующие функции по порядку
		fmt.Println(acceptMessage(m))
	}
}
