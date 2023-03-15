package main

import "fmt"

type system1 struct {
	name string
	id   int
}

func (s *system1) call() {
	fmt.Printf("called %s, id: %d\n", s.name, s.id)
}

type system2 struct {
	name string
	id   int
}

func (s *system2) call() {
	fmt.Printf("called %s, id: %d\n", s.name, s.id)
}

type system3 struct {
	name string
	id   int
}

func (s *system3) call() {
	fmt.Printf("called %s, id: %d\n", s.name, s.id)
}

type Facade struct {
	s1 system1
	s2 system2
	s3 system3
}

func GetFacade() *Facade {
	var f Facade
	f.s1 = system1{"system1", 10}
	f.s2 = system2{"system2", 20}
	f.s3 = system3{"system3", 30}
	return &f
}

// фасад объединяет внутренние вызовы функций в один простой вызов для клиента
func (f *Facade) Call() {
	f.s1.call()
	f.s2.call()
	f.s3.call()
}

func main() {
	f := GetFacade()
	f.Call()
}
