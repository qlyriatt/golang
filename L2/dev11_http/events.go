// events.go - функционал календаря
//
// функции работы с событиями требуют в качестве аргумента url.Values -
// переданные параметры и их значения
//
// query == http.Request.Form - для получения событий (метод GET)
// body == http.Request.PostForm - для получения событий (метод POST)
//
// функции отдают серверу []byte и ошибку, и не зависят от его внутреннего устройства

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"strconv"
)

// структура события
type Event struct {
	Name    string
	Hour    int
	Minutes int
	Day     int
	Month   int
	Year    int
}

// все события (кэш)
// функции редактируют эту переменную вместе с файлом events.json
var allEvents []Event

// получить события за день из общего списка
func eventsForDay(form url.Values) ([]byte, error) {

	// нужный день отображен в данном event
	searchFor, err := checkArgs(form)
	if err != nil {
		return []byte{}, err
	}

	// checkArgs() не проверяет данное условие
	// см. checkArgs()
	if searchFor.Day == 0 {
		return []byte{}, errors.New("отсутствует параметр - день")
	}

	events := make([]Event, 0)
	for _, event := range allEvents {

		if event.Day == searchFor.Day && event.Month == searchFor.Month && event.Year == searchFor.Year {
			events = append(events, event)
		}

	}

	// перевод событий в json в цикле, а не целым массивом
	// нужен для отображения каждого события на новой строке
	var data []byte
	for _, event := range events {
		e, err := json.Marshal(&event)
		// ошибка перевода в json
		if err != nil {
			return []byte{}, err
		}
		data = append(data, e...)
		data = append(data, []byte("\n")...)
	}

	return data, nil
}

// количество дней в месяце определенного года
func maxDaysInMonth(m int, y int) int {

	if m == 2 {
		if y%4 == 0 {
			return 29
		}

		return 28
	}

	if m == 4 || m == 6 || m == 9 || m == 11 {
		return 30
	}

	return 31
}

// получить события за неделю из общего списка
func eventsForWeek(form url.Values) ([]byte, error) {

	// дата начала недели отображена в данном event
	searchFor, err := checkArgs(form)
	if err != nil {
		return []byte{}, err
	}

	// checkArgs() не проверяет данное условие
	// см. checkArgs()
	if searchFor.Day == 0 {
		return []byte{}, errors.New("отсутствует параметр - день")
	}

	events := make([]Event, 0)
	d, m, y := searchFor.Day, searchFor.Month, searchFor.Year
	// вспомогательная структура, сопоставляет дни с месяцем для запрошенных событий
	months := make(map[int]int, 7)
	// вычисление дней для запрошенной недели
	// 7 дней
	for i := 0; i < 7; i++ {

		// число больше количества дней в данном месяце
		if d > maxDaysInMonth(m, y) {

			// переводим число в начало
			d = 1
			// если месяц последний, то переводим в начало, увеличиваем год
			if m == 12 {
				m = 1
				y++
			} else {
				// переключаем месяц на следующий
				m++
			}
		}

		// map[день] = месяц
		months[d] = m
		// переключаем день на следующий
		d++
	}

	// будет ли переключаться год
	yearInc := searchFor.Day >= 26 && searchFor.Month == 12
	// запрошенный год
	year := searchFor.Year
	for _, event := range allEvents {

		// если год переключался
		if yearInc {
			if event.Day >= 26 {
				// последние числа последнего месяца соответствуют запрошенному году
				year = searchFor.Year
			} else {
				// первые числа первого месяца соответствуют следующему году
				year = searchFor.Year + 1
			}
		}

		// если дня event.Day не существует в months, то months[event.Day] == 0,
		// условие не выполняется
		if months[event.Day] == event.Month && event.Year == year {
			events = append(events, event)
		}
	}

	// перевод событий в json в цикле, а не целым массивом
	// нужен для отображения каждого события на новой строке
	var data []byte
	for _, event := range events {
		e, err := json.Marshal(&event)
		// ошибка перевода в json
		if err != nil {
			return []byte{}, err
		}
		data = append(data, e...)
		data = append(data, []byte("\n")...)
	}

	return data, nil
}

func eventsForMonth(form url.Values) ([]byte, error) {

	// нужный месяц отображен в данном event
	searchFor, err := checkArgs(form)
	if err != nil {
		return []byte{}, err
	}

	// для правильного заполнения запроса
	//
	// checkArgs() не проверяет данное условие
	// см. checkArgs()
	if searchFor.Day != 0 {
		return []byte{}, errors.New("указан день в параметрах метода для получения информации о месяце")
	}

	events := make([]Event, 0)
	for _, event := range allEvents {

		if event.Month == searchFor.Month && event.Year == searchFor.Year {
			events = append(events, event)
		}
	}

	// перевод событий в json в цикле, а не целым массивом
	// нужен для отображения каждого события на новой строке
	var data []byte
	for _, event := range events {
		e, err := json.Marshal(&event)
		// ошибка перевода в json
		if err != nil {
			return []byte{}, err
		}
		data = append(data, e...)
		data = append(data, []byte("\n")...)
	}

	return data, nil
}

// существует ли данное событие
func eventExists(e Event) (int, bool) {
	for num, event := range allEvents {
		// ! прямое сравнение e == event может нестабильно работать с пустыми полями
		if e == event {
			return num, true
		}
	}

	return 0, false
}

// добавить событие в календарь
func addEvent(form url.Values) ([]byte, error) {

	// нужное событие отображено в eventToAdd
	eventToAdd, err := checkArgs(form)
	if err != nil {
		return []byte{}, err
	}

	// checkArgs() не проверяет данное условие
	// см. checkArgs()
	if eventToAdd.Day == 0 {
		return []byte{}, errors.New("отсутствует параметр - день")
	}

	if _, ok := eventExists(eventToAdd); ok {
		return []byte{}, errors.New("данное событие уже присутсвует в календаре")
	}

	data, err := json.Marshal(eventToAdd)
	// ошибка перевода в json
	if err != nil {
		return []byte{}, err
	}
	// добавление перехода на новую строку для ровного отображения вывода
	data = append(data, []byte("\n")...)

	file, err := os.OpenFile("events.json", os.O_APPEND|os.O_WRONLY, 0666)
	// ошибка открытия файла
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	// если произошла ошибка во время записи события в файл,
	// то оно не будет добавлено в allEvents, оставляя файл и кэш синхронизированными

	_, err = file.Write(data)
	if err != nil {
		return []byte{}, err
	}

	allEvents = append(allEvents, eventToAdd)

	return data, nil
}

// обновить событие
// удаляет событие и заменяет его на нужное
func updateEvent(form url.Values) ([]byte, error) {

	// отдельные url.Values для каждой операции (удаление + добавление)
	// см. checkArgs
	deleteForm := make(map[string][]string, len(form))
	addForm := make(map[string][]string, len(form))
	for key, vals := range form {
		if len(vals) < 2 {
			return []byte{}, errors.New("недостаточное количество значений для параметров")
		}

		// первое значение для параметра относится к событию для удаления
		deleteForm[key] = []string{vals[0]}
		// второе значение относится к событию для добавления на его место
		addForm[key] = []string{vals[1]}
	}

	_, err := deleteEvent(deleteForm)
	// ошибка в удалении
	if err != nil {
		return []byte{}, err
	}

	data, err := addEvent(addForm)
	// ошибка в добавлении
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// удалить событие из календаря
//
// ! в данный момент удаляет событие только из кэша для простоты тестирования
func deleteEvent(form url.Values) ([]byte, error) {

	// нужное событие отображено в eventToDelete
	eventToDelete, err := checkArgs(form)
	if err != nil {
		return []byte{}, err
	}

	// checkArgs() не проверяет данное условие
	// см. checkArgs()
	if eventToDelete.Day == 0 {
		return []byte{}, errors.New("отсутствует параметр - день")
	}

	// если событие существует
	if n, ok := eventExists(eventToDelete); ok {

		// все события после нужного смещаются на 1 к началу
		for i := n; i < len(allEvents)-1; i++ {
			allEvents[i] = allEvents[i+1]
		}
		// массив сокращается на 1 элемент
		allEvents = allEvents[:len(allEvents)-1]

	} else {
		return []byte{}, errors.New("событие отсутствует в календаре")
	}

	return []byte{}, nil
}

// парсинг url.Values для использования функциями
//
// получает url.Values - пары параметр-значение и преобразует их в Event
// каждая функция обрабатывает Event так, как ей нужно
func checkArgs(form url.Values) (Event, error) {

	var event Event
	for param, values := range form {

		// каждому параметру должно соответствовать одно значение
		// если функция должна получить несколько значений для одного параметра,
		// она должна вызвать checkArgs() несколько раз с подготовленными url.Values
		// см. updateEvent()
		if len(values) > 1 {
			return Event{}, errors.New("параметры со множественными значениями не поддерживается")
		}

		if param == "name" {
			event.Name = values[0]
			continue
		}

		// если параметр не "name", то его значение может быть только int, записывается в val
		val, err := strconv.Atoi(values[0])
		if err != nil {
			return Event{}, errors.New("ошибка в преобразовании значения параметра/неподдерживаемый параметр")
		}

		switch param {

		case "year":

			if val < 1970 || val > 2050 {
				return Event{}, errors.New("неверный год")
			}
			event.Year = val

		case "month":

			if val < 1 || val > 12 {
				return Event{}, errors.New("неверный месяц")
			}
			event.Month = val

		case "day":

			// первая проверка дня
			if val < 1 || val > 31 {
				return Event{}, errors.New("неверный день")
			}
			event.Day = val

		case "hour":

			if val < 0 || val > 23 {
				return Event{}, errors.New("неверное значение часа")
			}
			event.Hour = val

		case "minutes":

			if val < 0 || val > 59 {
				return Event{}, errors.New("неверное значение минут")
			}
			event.Minutes = val

		default:
			return Event{}, errors.New("встречен неподдерживаемый параметр")
		}
	}

	d := event.Day
	m := event.Month
	y := event.Year

	// не проверяет d == 0, поскольку такой параметр нужно передать в GET -> events_for_month
	// см. eventsForMonth()
	//
	// функции, которым нужна эта проверка, должны осуществить её самостоятельно
	// например, см. addEvent()
	if y == 0 || m == 0 /* || d == 0 */ {
		return Event{}, errors.New("недостаточно параметров для определения события")
	}

	// дальнейшие проверки дня

	if m == 2 {
		if d == 29 && (y%4 != 0) {
			return Event{}, errors.New("29 число в феврале невисокосного года")
		} else if d > 29 {
			return Event{}, errors.New("неподходящее число для февраля")
		}
	}

	if d == 31 && (m == 4 || m == 6 || m == 9 || m == 11) {
		return Event{}, errors.New("31 число в неподходящем месяце")
	}

	return event, nil
}

// создание случайного списка событий
func createEvents() error {

	const eventsToCreate = 100

	// событие получает случайно
	// день 1-31
	// месяц 1-12
	// год 2000-2019
	for i, numEvents := 0, eventsToCreate; i < numEvents; i++ {

		y := 2000 + rand.Intn(20)
		m := 1 + rand.Intn(12)
		d := 1 + rand.Intn(30)
		if m == 2 && d >= 28 {
			d = 1 + rand.Intn(28)
		}

		var event Event
		event.Day = d
		event.Month = m
		event.Year = y

		allEvents = append(allEvents, event)
	}

	// события сортируются по возрастанию в порядке год -> месяц -> день
	sort.Slice(allEvents, func(i, j int) bool {
		if allEvents[j].Year > allEvents[i].Year {
			return true
		} else if allEvents[j].Year < allEvents[i].Year {
			return false
		}

		if allEvents[j].Month > allEvents[i].Month {
			return true
		} else if allEvents[j].Month < allEvents[i].Month {
			return false
		}

		if allEvents[j].Day > allEvents[i].Day {
			return true
		} else if allEvents[j].Day < allEvents[i].Day {
			return false
		}

		return true
	})

	file, err := os.Create("events.json")
	// ошибка в создании файла
	if err != nil {
		return err
	}
	defer file.Close()

	// после этого момента события присутствуют в allEvents, файл events.json создан
	//
	// если дальше происходит ошибка, в файле окажутся некорректные/неполные данные,
	// а простая проверка на присутствие файла этого не отразит
	//
	// т.к. при таком варианте в allEvents присутствуют события, достаточно проверить len(allEvents) == 0,
	// и в случае неравенства запустить функцию создания еще раз
	// см. loadEvents()

	// перевод в json
	var data []byte
	for _, event := range allEvents {
		e, err := json.Marshal(event)
		// ошибка в переводе в json
		if err != nil {
			return err
		}
		data = append(data, e...)
		data = append(data, []byte("\n")...)
	}

	_, err = file.Write(data)
	// ошибка в записи в файл
	if err != nil {
		return err
	}

	return nil
}

// загрузка событий из файла в кэш
func loadEvents() error {

	// файл существует и в кэше нет записей
	if _, err := os.Stat("events.json"); err != os.ErrNotExist && len(allEvents) == 0 {

		file, err := os.Open("events.json")
		// ошибка открытия файла
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// ошибка сканирования
			if scanner.Err() != nil {
				return scanner.Err()
			}

			var event Event
			err := json.Unmarshal([]byte(scanner.Text()), &event)
			// ошибка перевода json в Event
			if err != nil {
				return err
			}
			allEvents = append(allEvents, event)
		}

	} else {
		// если файла не существует
		// или в кэше уже присутствуют записи при существующем файле и попытке запустить загрузку ( см. createEvents() )
		// необходимо запустить функцию создания вместо загрузки
		err := createEvents()
		if err != nil {
			return err
		}
	}

	return nil
}
