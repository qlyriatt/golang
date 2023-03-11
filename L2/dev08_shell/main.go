package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-ps"
)

// вывод текущей директории
func getColoredDir() string {
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	wd = strings.Replace(wd, home, "~", 1)
	return "\033[1;96m" + wd + "\033[0m" + "$ "
}

type command struct {
	cmd  string   // команда
	args []string // аргументы
	pipe bool     // отдает ли команда вывод в pipe
}

// доступные команды
var commands = map[string]bool{
	"cd":   true,
	"pwd":  true,
	"echo": true,
	"kill": true,
	"ps":   true,
	"exec": true,
	"fork": true,
	"exit": true,
}

func parseInput(input string) (cmds []command, err error) {

	fields := strings.Fields(input) // все поля входной строки
	thisFieldNum := 0               // номер поля строки fields для текущей команды
	thisCommand := command{}        // буфер текущей команды
	for fieldNum := 0; fieldNum < len(fields); fieldNum++ {

		// первое поле для каждой thisCommand - сама команда
		if thisFieldNum == 0 {
			if commands[fields[fieldNum]] {
				thisCommand.cmd = fields[fieldNum]
				thisFieldNum++
				continue
			} else {
				// команды не существует
				return cmds, errors.New(fields[fieldNum] + ": is not a valid command")
			}
		}

		// разделители команд
		// добавление готовой команды в массив команд cmds,
		// переход к конструированию следующей команды
		if fields[fieldNum] == "|" || fields[fieldNum] == "&&" {
			// pipe
			if fields[fieldNum] == "|" {
				thisCommand.pipe = true
			}

			// && - выполнение команд подряд
			cmds = append(cmds, thisCommand)
			thisCommand = command{}
			thisFieldNum = 0
			continue
		}

		// добавить аргумент для данной команды, перейти к следующему полю
		thisCommand.args = append(thisCommand.args, fields[fieldNum])
		thisFieldNum++
	}
	// добавление последней готовой команды в cmds
	cmds = append(cmds, thisCommand)
	return cmds, nil
}

func cd(args []string) error {

	if len(args) > 1 {
		return errors.New("cd: слишком много аргументов")
	}
	if len(args) < 1 {
		return nil
	}

	// смена директории
	if err := os.Chdir(args[0]); err != nil {
		return err
	}

	return nil
}

func pwd() (string, error) {

	// получение текущей директории
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil

}

func echo(args []string) (string, error) {

	// сбор общей строки из всех args
	var buf string
	for _, arg := range args {
		buf += arg + " "
	}
	return buf, nil
}

func kill(args []string) error {

	if len(args) < 1 {
		return errors.New("kill: недостаточно аргументов")
	}

	// проверяем все args
	for _, arg := range args {
		// данный arg - не int
		if pid, err := strconv.Atoi(arg); err != nil {
			fmt.Println(err)
			continue
		} else if proc, err := ps.FindProcess(pid); err != nil {
			// ошибка в поиске процесса
			fmt.Println(err)
			continue
		} else if proc == nil {
			// процесса с таким pid нет
			fmt.Println("kill:", pid, ": процесс не существует")
			continue
		} else if proc, err := os.FindProcess(pid); err != nil {
			// ошибка в поиске процесса
			fmt.Println(err)
		} else {
			// отправка сигнала interrupt процессу
			proc.Signal(os.Interrupt)
		}
	}

	return nil
}

func psn() (procs []string, err error) {

	ps, err := ps.Processes()
	if err != nil {
		return []string{}, err
	}

	for _, p := range ps {
		procs = append(procs, fmt.Sprintf("%7v %v\n", p.Pid(), p.Executable()))
	}
	// возвращает массив, где каждая строка это процесс
	return procs, nil

}

func ex(args []string) error {

	if len(args) < 1 {
		return errors.New("exec: недостаточно аргументов")
	}

	argv := []string{}
	if len(args) > 1 {
		argv = args[1:]
	}

	// вызов exec для args[0]
	if err := syscall.Exec(args[0], argv, os.Environ()); err != nil {
		return err
	}

	return nil
}

func fork() error {

	// поиск программы shell (qsh)
	path, err := os.Executable()
	if err != nil {
		return err
	}

	// запуск qsh
	cmd := exec.Command(path)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil

}

func main() {

	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println()
	for {
		// вывод пригласительной строки
		fmt.Print(getColoredDir())

		if scanner.Scan(); scanner.Err() != nil {
			fmt.Println(scanner.Err())
			return
		}
		input = scanner.Text()

		cmds, err := parseInput(input)
		if err != nil {
			fmt.Println(err)
		}

		pipedArgs := []string{}
		piped := false
		for _, cmd := range cmds {

			// если предыдущая команда отправила вывод в pipe,
			// то текущая команда получает его как аргументы
			if piped {
				cmd.args = pipedArgs
				pipedArgs = []string{}
				piped = false
			}

			var out = []string{} // вывод текущей команды, может быть пустым

			// если внутри одной из команд в цепочке (&& или |) происходит ошибка,
			// то цепочка прерывается, отображается результат только
			// успешно выполнившихся команд
			if cmd.cmd == "cd" {
				if err := cd(cmd.args); err != nil {
					fmt.Println(err)
					break
				}
			} else if cmd.cmd == "pwd" {
				if res, err := pwd(); err != nil {
					fmt.Println(err)
					break
				} else {
					out = append(out, res)
				}
			} else if cmd.cmd == "echo" {
				if res, err := echo(cmd.args); err != nil {
					fmt.Println(err)
					break
				} else {
					out = append(out, res)
				}
			} else if cmd.cmd == "kill" {
				if err := kill(cmd.args); err != nil {
					fmt.Println(err)
					break
				}
			} else if cmd.cmd == "ps" {
				if res, err := psn(); err != nil {
					fmt.Println(err)
					break
				} else {
					out = append(out, res...)
				}
			} else if cmd.cmd == "exec" {
				if err := ex(cmd.args); err != nil {
					fmt.Println(err)
					break
				} else {
					return
				}
			} else if cmd.cmd == "fork" {
				go fork()
			} else if cmd.cmd == "exit" {
				fmt.Println()
				return
			}

			// если команда передает свой результат в pipe, добавить его в pipedArgs
			if cmd.pipe {
				pipedArgs = append(pipedArgs, out...)
				piped = true
			} else {
				// иначе вывести на экран
				for _, str := range out {
					fmt.Print(str)
				}
				fmt.Println()
			}
		}
	}
}
