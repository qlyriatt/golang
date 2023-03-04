package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	// -fstartOn
	// -fstartOn-
	// -fstartOn-endOn
	// -f-endOn

	startOn := 1
	endOn := 0
	delim := "\t"
	skipNoDelim := false

	for _, arg := range os.Args[1:] {

		if strings.Contains(arg, "-f") {
			arg := arg[2:]

			if len(arg) < 1 {
				fmt.Println("err: -f: no fields selected")
				return
			}

			// -f-
			if arg[0] == '-' && len(arg) == 1 {
				fmt.Println("err: -f: wrong field selection")
				return
			}

			if arg[0] == '-' {
				var err error
				if endOn, err = strconv.Atoi(arg[1:]); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				isRange := false
				for i, ch := range arg {
					if ch == '-' {

						isRange = true
						var err error

						// startOn
						if i != 0 {
							if startOn, err = strconv.Atoi(arg[0:i]); err != nil {
								fmt.Println(err)
								return
							}
							if startOn == 0 {
								fmt.Println("err: -f: fields start at 1")
								return
							}
						}

						// endOn
						if i != len(arg)-1 {
							if endOn, err = strconv.Atoi(arg[i+1:]); err != nil {
								fmt.Println(err)
								return
							}
							if endOn < startOn {
								fmt.Println("err: -f: wrong range")
								return
							}
						}

					}
				}
				if !isRange {
					var err error
					if startOn, err = strconv.Atoi(arg); err != nil {
						fmt.Println(err)
						return
					}
					endOn = -1
				}
			}

		} else if strings.Contains(arg, "-d") {
			arg := arg[2:]
			log.Printf("arg: %q", arg)
			// if len(arg) < 1 {
			// 	fmt.Println("err: -d: no delimeter selected")
			// 	return
			// }

			// first, last := -1, 0
			// for i, ch := range arg {
			// 	if ch == '\'' {
			// 		if first == -1 {
			// 			first = i
			// 		} else {
			// 			last = i
			// 		}
			// 	}
			// }
			// if first == -1 || last == 0 {
			// 	fmt.Println("err: -d: error in delimeter quotes")
			// 	return
			// }

			delim = arg
		} else if strings.Contains(arg, "-s") && len(arg) == 2 {
			skipNoDelim = true
		} else {
			fmt.Println("err: wrong arguments")
			return
		}
	}

	log.Println("start:", startOn, "end:", endOn)

	file, _ := os.Open("test.txt")
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(os.Stdin)

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			return
		}

		str := scanner.Text()

		// -s
		if skipNoDelim && !strings.Contains(str, delim) {
			continue
		}

		fields := strings.SplitAfter(scanner.Text(), delim)
		for i := range fields {
			fields[i], _ = strings.CutSuffix(fields[i], delim)
		}

		if len(fields) < startOn {
			continue
		}

		if endOn == -1 || startOn == endOn { // одно поле
			fields = fields[startOn-1 : startOn]
		} else if endOn == 0 {
			fields = fields[startOn-1:]
		} else if endOn > len(fields) {
			fields = fields[startOn-1:]
		} else {
			fields = fields[startOn-1 : endOn]
		}

		for _, field := range fields {
			if _, err := writer.WriteString(field + " \n"); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(err)
		return
	}
}
