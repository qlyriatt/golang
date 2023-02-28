package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type flags struct {
	after      int
	before     int
	context    int  // приоритет над after, before
	count      bool // приоритет над after, before, context, lineNum
	ignoreCase bool
	invert     bool // приоритет над after, before, context
	fixed      bool
	lineNum    bool
}

func output(writer *bufio.Writer, str string, i int, flags flags) error {

	if flags.lineNum { // -n
		if _, err := writer.WriteString(strconv.Itoa(i) + ": "); err != nil {
			return err
		}
	}

	if _, err := writer.WriteString(str + "\n"); err != nil {
		return err
	}

	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("err: no filename")
		return
	}

	var filename string
	if _, err := os.Stat(args[len(args)-1]); errors.Is(err, os.ErrNotExist) {
		fmt.Println("err:", args[len(args)-1], "- no such file")
		return
	} else if err != nil {
		fmt.Println(err)
		return
	} else if len(args) == 1 {
		fmt.Println("err: type something to search for")
		return
	} else {
		filename = args[len(args)-1]
		args = args[:len(args)-1]
	}

	var flags flags
	var pattern string
	var foundPattern = false

	// -A -B -C -c -i -v -F -n
	for _, arg := range args {

		if strings.Contains(arg, "-A") && len(arg) > 2 {
			if val, err := strconv.Atoi(arg[2:]); err != nil {
				fmt.Println(err)
				return
			} else if val <= 0 {
				fmt.Println("err: -A: invalid length argument")
				return
			} else {
				flags.after = val
			}
		} else if strings.Contains(arg, "-B") && len(arg) > 2 {
			if val, err := strconv.Atoi(arg[2:]); err != nil {
				fmt.Println(err)
				return
			} else if val <= 0 {
				fmt.Println("err: -B: invalid length argument")
				return
			} else {
				flags.before = val
			}
		} else if strings.Contains(arg, "-C") && len(arg) > 2 {
			if val, err := strconv.Atoi(arg[2:]); err != nil {
				fmt.Println(err)
				return
			} else if val <= 0 {
				fmt.Println("err: -C: invalid length argument")
				return
			} else {
				flags.context = val
			}
		} else if strings.Contains(arg, "-c") && len(arg) == 2 {
			flags.count = true
		} else if strings.Contains(arg, "-i") && len(arg) == 2 {
			flags.ignoreCase = true
		} else if strings.Contains(arg, "-v") && len(arg) == 2 {
			flags.invert = true
		} else if strings.Contains(arg, "-F") && len(arg) == 2 {
			flags.fixed = true
		} else if strings.Contains(arg, "-n") && len(arg) == 2 {
			flags.lineNum = true
		} else if !foundPattern {
			foundPattern = true
			pattern = arg
		}
	}

	input, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer input.Close()

	writeAfter := 0
	writeBefore := 0
	if flags.context != 0 {
		writeAfter = flags.context
		writeBefore = flags.context
	} else {
		if flags.after != 0 {
			writeAfter = flags.after
		}

		if flags.before != 0 {
			writeBefore = flags.before
		}
	}

	bufBefore := make([]string, 0, writeBefore)
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(os.Stdout)
	count := 0                           // количество нужных строк
	i := -1                              // номер текущей строки
	printed := make(map[int]struct{}, 0) // номера напечатанных строк
	leftToWriteAfter := 0                // осталось напечатать любых строк после подходящей
	for scanner.Scan() {

		// ошибка сканнера
		if scanner.Err() != nil {
			fmt.Println(err)
			return
		}

		i++
		str := scanner.Text()

		// -i
		if flags.ignoreCase {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}

		// строка не содержит паттерн
		if !strings.Contains(str, pattern) {

			if !flags.invert {

				// -B
				// буфер предыдущих строк
				if writeBefore != 0 {

					// заполнение буфера предыдущих строк
					if len(bufBefore) != writeBefore {
						bufBefore = append(bufBefore, str)
					} else { // при заполненности - прокрутка(цикл) буфера предыдущих строк

						for i := 0; i < len(bufBefore)-1; i++ {
							bufBefore[i] = bufBefore[i+1]
						}
						bufBefore[len(bufBefore)-1] = str
					}
				}

				// -A
				// печать последующих строк
				if leftToWriteAfter > 0 {
					output(writer, str, i, flags)
					printed[i] = struct{}{}
					leftToWriteAfter--
				}
			} else if flags.count { // -v -c
				count++
			} else { // -v

				if err := output(writer, str, i, flags); err != nil {
					fmt.Println(err)
					return
				}
				printed[i] = struct{}{}
			}
			continue
		}

		// -F
		// строка содержит нужный паттерн, но не совпадает точно
		if flags.fixed && str != pattern {
			continue
		}

		// -c
		if flags.count {
			count++
			continue
		}

		// -v
		if flags.invert {
			continue
		}

		// -B/-C
		if writeBefore != 0 {
			for b, bufStr := range bufBefore {

				// строка в буфере уже была напечатана
				if _, ok := printed[i-(len(bufBefore)-b)]; ok {
					continue
				}

				if err := output(writer, bufStr, i-(len(bufBefore)-b), flags); err != nil {
					fmt.Println(err)
					return
				}
				printed[i-(len(bufBefore)-b)] = struct{}{}
			}

			if err := writer.Flush(); err != nil {
				fmt.Println(err)
				return
			}
		}

		// печать подходящей строки
		if err := output(writer, str, i, flags); err != nil {
			fmt.Println(err)
			return
		}
		printed[i] = struct{}{}

		if writeAfter != 0 {
			leftToWriteAfter = writeAfter
		}
	}

	if flags.count {
		fmt.Println(count)
	} else {
		if err := writer.Flush(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
