package main

import "fmt"

// интерфейс команды
type command interface {
	execute()
	undo()
}

// получатель команд
type document struct {
	text []string
}

func (d *document) insert(s string) {
	d.text = append(d.text, s)
}

// конкретная имплементация команды для вызывающего объекта - insert
type insertCommand struct {
	s string
	d *document
}

func (i *insertCommand) execute() {
	i.d.insert(i.s)
}

func (i *insertCommand) undo() {
	i.d.text = i.d.text[:len(i.d.text)-1]
}

// вызывающий команды объект
type viewer struct {
	commands []command // буфер выполненных команд
	d        *document //
}

// вставка
func (v *viewer) Insert(s string) {
	v.d.insert(s)
	v.commands = append(v.commands, &insertCommand{s, v.d})
}

// отмена
func (v *viewer) Undo() {

	if len(v.commands) <= 0 {
		return
	}

	v.commands[len(v.commands)-1].undo()
	v.commands = v.commands[:len(v.commands)-1]
}

func main() {
	doc := document{[]string{"hello", "world", "hello", "world"}}
	v := viewer{[]command{}, &doc}

	// не вызываются команды для работы над document
	//
	// вместо этого работаем с вызывающим объектом - viewer
	v.Insert("hi")
	v.Insert("hi again")
	v.Undo()

	fmt.Println(v.d)
}
