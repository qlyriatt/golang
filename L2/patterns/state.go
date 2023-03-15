package main

import "fmt"

// интерфейс структур состояния
type state interface {
	onClick()
	stateName()
	hide(w *widget)
	show(w *widget)
}

type widget struct {
	state
}

func (w *widget) Click() {
	w.state.onClick()
}

func (w *widget) WhatState() {
	w.state.stateName()
}

func (w *widget) Hide() {
	w.state.hide(w)
}

func (w *widget) Show() {
	w.state.show(w)
}

// скрытое состояние виджета
type hidden struct {
}

func (h hidden) onClick() {
	fmt.Println("nothing happened on click")
}

func (h hidden) stateName() {
	fmt.Println("hidden widget")
}

func (h hidden) hide(w *widget) {
	fmt.Println("already hidden")
}

func (h hidden) show(w *widget) {
	w.state = visible{}
	fmt.Println("widget is visible now")
}

// видимое состояние виджета
type visible struct {
}

func (v visible) onClick() {
	fmt.Println("clicked widget")
}

func (v visible) stateName() {
	fmt.Println("visible widget")
}

func (v visible) hide(w *widget) {
	w.state = hidden{}
	fmt.Println("widget is hidden now")
}

func (v visible) show(w *widget) {
	fmt.Println("already visible")
}

func main() {

	var w = widget{visible{}}
	// в зависимости от текущего состояния (поля state)
	// widget по разному реагирует на внешние запросы

	// visible
	w.Click()
	w.WhatState()
	w.Show()
	w.Hide()

	fmt.Println()

	// hidden
	w.Click()
	w.WhatState()
	w.Hide()
	w.Show()
}
