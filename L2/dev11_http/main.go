package main

import (
	"fmt"
)

func main() {

	// загрузка событий
	err := loadEvents()
	if err != nil {
		fmt.Println(err)
		return
	}

	// запуск сервера
	err = startServer()
	if err != nil {
		fmt.Println(err)
		return
	}
}
