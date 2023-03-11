// server.go - функционал сервера
//
// с устройством календаря сервер работает только по определенным предоставленным функциям -
// eventsForDay(), addEvent() и аналоги
//
// сервер не содержит в себе информации о том, как внутренне устроена логика работы календаря,
// например, не знает о типе Event и переменной allEvents

package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"
)

// добавить логирование к handler функции
func addLogs(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// попытка перенаправить вывод в логи
		output := os.Stdout
		file, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Fprintln(output, "ошибка открытия файла логов, вывод в stdout")
		} else {
			output = file
		}

		// метод и url запроса
		fmt.Fprintf(output, "%v\nRequest: %s -> %s\n", time.Now().Format(time.TimeOnly), r.Method, r.URL)

		r.ParseForm()
		// полученные параметры query для метода GET
		if r.Method == http.MethodGet {
			fmt.Fprint(output, "Query: ")
			for arg, val := range r.Form {
				fmt.Fprintf(output, "%s: %s ", arg, val)
			}
		}

		// полученные параметры body для метода POST
		if r.Method == http.MethodPost {
			fmt.Fprint(output, "Body: ")
			for arg, val := range r.PostForm {
				fmt.Fprintf(output, "%s: %s ", arg, val)
			}
		}

		// вместо настоящего вывода, ответ функции сначала записывается в Recorder
		resp := httptest.NewRecorder()

		// непосредственно handler функция, отвечает на запрос
		h(resp, r)

		// ответ
		fmt.Fprintf(output, "\nResponse: %v\n", resp.Result().StatusCode)
		// убрать последний \n перед записью в логи, чтобы не оставлять пустую строку
		if len(resp.Body.Bytes()) > 2 {
			removeNewline := string(resp.Body.Bytes())[:len(resp.Body.Bytes())-1]
			fmt.Fprintf(output, "%v\n", removeNewline)
		}
		// разделитель запросов
		fmt.Fprintf(output, "%v\n", strings.Repeat("=", 70))

		// копирование ответа функции с Recorder на настоящий вывод
		for k, vals := range resp.Result().Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.Result().StatusCode)
		w.Write(resp.Body.Bytes())
	}
}

// получить handler для GET запроса
func getHandler(url string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// если метод - GET
		if r.Method == http.MethodGet {

			err := r.ParseForm()
			// ошибка в парсинге параметров
			if err != nil {
				w.WriteHeader(500)
				return
			}

			// ответ от функции календаря записывается в data
			var data []byte
			switch url {
			case "events_for_day":
				data, err = eventsForDay(r.Form)
			case "events_for_week":
				data, err = eventsForWeek(r.Form)
			case "events_for_month":
				data, err = eventsForMonth(r.Form)
			default:
				w.WriteHeader(500)
				return
			}

			// функция из switch вернула ошибку в параметрах
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				// добавление перехода на новую строку для ровного отображения вывода
				w.Write([]byte("\n"))
				return
			}

			// если подходящих событий не нашлось, то массив data будет пуст
			if len(data) == 0 {
				w.WriteHeader(204)
				return
			}

			w.Write(data)

		} else {
			// неподдерживаемый метод
			w.WriteHeader(405)
		}
	}
}

// получить handler для POST запроса
func postHandler(url string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// если метод - POST
		if r.Method == http.MethodPost {

			err := r.ParseForm()
			// ошибка в парсинге параметров
			if err != nil {
				w.WriteHeader(500)
				return
			}

			// ответ от функции календаря записывается в data
			var data []byte
			switch url {
			case "create_event":
				data, err = addEvent(r.PostForm)
			case "update_event":
				data, err = updateEvent(r.PostForm)
			case "delete_event":
				data, err = deleteEvent(r.PostForm)
			default:
				w.WriteHeader(500)
				return
			}

			// функция из switch вернула ошибку в параметрах
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				// добавление перехода на новую строку для ровного отображения вывода
				w.Write([]byte("\n"))
				return
			}

			// пустой []byte без ошибки возвращается в ответ на case "delete_event"
			if len(data) == 0 {
				w.WriteHeader(204)
				return
			}

			w.Write(data)

		} else {
			// неподдерживаемый метод
			w.WriteHeader(405)
		}
	}
}

// запустить сервер
func startServer() error {

	data, err := os.ReadFile("config.txt")
	// отсутствует конфигурационный файл
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(string(data))
	// ошибка преобразования для номера из файла
	if err != nil {
		return err
	}
	if port < 1000 || port > (1<<16) {
		return errors.New("недопустимое значение порта в config.txt")
	}

	// регистрация функций в стандартном роутере (http.DefaultServeMux)
	http.Handle("/events_for_day", addLogs(getHandler("events_for_day")))
	http.Handle("/events_for_week", addLogs(getHandler("events_for_week")))
	http.Handle("/events_for_month", addLogs(getHandler("events_for_month")))
	http.Handle("/create_event", addLogs(postHandler("create_event")))
	http.Handle("/update_event", addLogs(postHandler("update_event")))
	http.Handle("/delete_event", addLogs(postHandler("delete_event")))

	// запуск сервера на выбранном порту
	err = http.ListenAndServe("localhost:"+strconv.Itoa(port), nil)
	if err != nil {
		return err
	}

	return nil
}
