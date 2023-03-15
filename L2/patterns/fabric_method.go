package main

import "fmt"

// интерфейс для всех painter
type painter interface {
	chooseColor() // своя реализация для каждого painter
	mixPaint()    // общая для всех painter реализация
	paint()       // общая для всех painter реализация
}

// "базовый" класс, не умеет выбирать цвет
// определяет общие функции mixPaint() и paint()
type basicPainter struct {
	color string
}

func (p basicPainter) mixPaint() {
	fmt.Println("mixed", p.color, "paint")
}

func (p basicPainter) paint() {
	fmt.Println("painted with", p.color, "color")
}

// определяет собственную реализацию chooseColor()
type bluePainter struct {
	basicPainter
}

func (p *bluePainter) chooseColor() {
	p.color = "blue"
	fmt.Println("chosen blue color")
}

// определяет собственную реализацию chooseColor()
type greenPainter struct {
	basicPainter
}

func (p *greenPainter) chooseColor() {
	p.color = "green"
	fmt.Println("chosen green color")
}

func main() {
	painters := []painter{&greenPainter{}, &bluePainter{}}

	for _, p := range painters {
		p.chooseColor()
		p.mixPaint()
		p.paint()
		fmt.Println()
	}
}
