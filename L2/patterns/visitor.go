package main

import "fmt"

// отдельный визитор определен для каждой системы (класса)
type visitor interface {
	visitString(stringSystem)
	visitInt(intSystem)
	visitBool(boolSystem)
}

// каждая система принимает визитора
type system interface {
	accept(visitor)
}

type stringSystem struct {
	field string
}

func (s stringSystem) accept(v visitor) {
	v.visitString(s)
}

type intSystem struct {
	field int
}

func (i intSystem) accept(v visitor) {
	v.visitInt(i)
}

type boolSystem struct {
	field bool
}

func (b boolSystem) accept(v visitor) {
	v.visitBool(b)
}

// визитор, показывающий тип системы
//
// для каждой конкретной системы имеет свою реализацию
type TypeVisitor struct {
	t string
}

func (tv *TypeVisitor) visitString(s stringSystem) {
	tv.t = "string"
}

func (tv *TypeVisitor) visitInt(i intSystem) {
	tv.t = "int"
}

func (tv *TypeVisitor) visitBool(b boolSystem) {
	tv.t = "bool"
}

// визитор, показывающий значение системы
//
// для каждой конкретной системы имеет свою реализацию
type ValueVisitor struct {
	v any
}

func (vv *ValueVisitor) visitString(s stringSystem) {
	vv.v = s.field
}

func (vv *ValueVisitor) visitInt(i intSystem) {
	vv.v = i.field
}

func (vv *ValueVisitor) visitBool(b boolSystem) {
	vv.v = b.field
}

func main() {

	// системы различного вида
	var systems = []system{&stringSystem{"hello"}, &intSystem{22}, &boolSystem{true}}

	var tv = &TypeVisitor{}
	var vv = &ValueVisitor{}
	for _, s := range systems {
		// каждая система принимает визитора, определенного конкретно для нее
		s.accept(tv)
		s.accept(vv)
		fmt.Println(tv.t, vv.v)
	}
}
