package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// нет имени файла
	if len(os.Args) == 1 {
		fmt.Println("no filename")
		return
	}
	args := os.Args[1 : len(os.Args)-1]

	// невозможно открыть файл/неправильное имя файла
	input, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer input.Close()

	// флаги
	sortColumn := 1      // -k
	sortNums := false    // -n
	sortReverse := false // -r
	sortUnique := false  // -u
	for _, arg := range args {
		if strings.Contains(arg, "-k") {

			if len(arg) < 3 {
				fmt.Println("-k: should have a column number")
				return
			}

			if val, err := strconv.Atoi(arg[2:]); err != nil {
				fmt.Println("-k: ", err)
				return
			} else if val <= 0 {
				fmt.Println("-k: wrong column number")
				return
			} else {
				sortColumn = val
			}

		} else if strings.Contains(arg, "-n") && len(arg) == 2 {
			sortNums = true
		} else if strings.Contains(arg, "-r") && len(arg) == 2 {
			sortReverse = true
		} else if strings.Contains(arg, "-u") && len(arg) == 2 {
			sortUnique = true
		} else {
			fmt.Println("invalid arguments")
			return
		}
	}

	scanner := bufio.NewScanner(input)

	data := make([]string, 0)        // входные строки из файла
	unsorted := make(map[int]string) // колонка для сортировки
	str := -1
	for scanner.Scan() {
		if scanner.Err() != nil {
			fmt.Println(err)
			return
		}
		str++
		data = append(data, scanner.Text())
		fields := strings.Fields(scanner.Text())

		if sortColumn > len(fields)+1 {
			continue
		}

		unsorted[str] = fields[sortColumn-1]
	}

	keys := make([]int, 0)
	values := make(map[string]struct{})
	for k := range unsorted {
		if sortUnique {
			if _, ok := values[unsorted[k]]; ok {
				continue
			}
			values[unsorted[k]] = struct{}{}
		}
		keys = append(keys, k)
	}

	defer func() { // ловит panic() из less() в SliceStable()
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()
	sort.SliceStable(keys, func(i, j int) bool {
		var tmp []string
		tmp = append(tmp, unsorted[keys[i]])
		tmp = append(tmp, unsorted[keys[j]])

		if !sortNums {
			if int([]rune(tmp[0])[0]) > int([]rune(tmp[1])[0]) {
				tmp[0], tmp[1] = tmp[1], tmp[0]
			}
		} else {
			i1, err := strconv.Atoi(unsorted[keys[i]])
			if err != nil {
				panic(err)
			}
			i2, err := strconv.Atoi(unsorted[keys[j]])
			if err != nil {
				panic(err)
			}

			if i1 > i2 {
				tmp[0], tmp[1] = tmp[1], tmp[0]
			}
		}

		if !sortReverse {
			if tmp[0] == unsorted[keys[i]] {
				return true
			}
			return false

		}

		if tmp[0] == unsorted[keys[i]] {
			return false
		}
		return true

	})

	// вывод
	output, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	for _, k := range keys {
		if _, err := writer.WriteString(data[k] + "\n"); err != nil {
			fmt.Println(err)
			return
		}
		if err := writer.Flush(); err != nil {
			fmt.Println(err)
			return
		}
	}

}
