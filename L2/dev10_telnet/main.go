package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"os"
	"time"
)

var timeout int

func main() {

	flag.IntVar(&timeout, "timeout", 10, "таймаут подключения")
	flag.Parse()

	if len(flag.Args()) > 2 {
		fmt.Println("встречены лишние аргументы")
		return
	}

	if len(flag.Args()) < 2 {
		fmt.Println("недостаточно аргументов")
		return
	}

	var addr, port string
	ip, err := netip.ParseAddr(flag.Args()[0])
	// если не удалось привести первый аргумент к виду IP
	if err != nil {
		// пробуем сделать адрес доменным именем
		addr = flag.Args()[0]
	} else {
		// делаем адрес IP
		addr = ip.String()
	}
	// порт
	port = flag.Args()[1]

	var conn net.Conn   // соединение
	start := time.Now() // время начала попытки подключения
	// пока не истекло время, пробуем установить соединение
	for time.Now().Sub(start).Seconds() < float64(timeout) {
		var err error
		conn, err = net.Dial("tcp", addr+":"+port)
		// при неудаче, sleep(1) и повторная попытка
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	// не удалось установить соединение
	if conn == nil {
		fmt.Println("таймаут подключения")
		return
	}

	defer conn.Close()

	done := make(chan struct{})

	go func() {

		inputUser := bufio.NewScanner(os.Stdin)
		fmt.Print("out>")

		// сканирование ввода пользователя
		for inputUser.Scan() {

			if inputUser.Err() != nil {
				fmt.Println(err)
				break
			}

			fmt.Print("out>")

			_, err := conn.Write([]byte(inputUser.Text()))
			// ошибка записи
			if err != nil {
				fmt.Println(err)
				continue
			}

		}
		// конец цикла означает, что пользователь прервал ввод,
		// нужно закрыть соединение

		// Close() завершает цикл сканирования ввода с сокета
		conn.Close()
	}()

	go func() {

		inputSocket := bufio.NewScanner(conn)

		// сканирование ввода с сокета
		for inputSocket.Scan() {

			if inputSocket.Err() != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("\nin>", inputSocket.Text())

		}
		// цикл завершается при закрытии соединения

		done <- struct{}{}
	}()

	// ожидание завершения рутин
	<-done
	fmt.Println()
}
