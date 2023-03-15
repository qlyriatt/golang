package main

import "fmt"

// общий интерфейс builder
type houseBuilder interface {
	buildRoof()
	buildWalls()
	buildFoundation()
	GetHouse() *house
}

// структура для builder
type house struct {
	roof,
	walls,
	foundation string
}

// конкретный объект
type modernHouseBuilder struct {
	h house
}

func (b *modernHouseBuilder) buildRoof() {
	b.h.roof = "glass"
}

func (b *modernHouseBuilder) buildWalls() {
	b.h.walls = "brick"
}

func (b *modernHouseBuilder) buildFoundation() {
	b.h.foundation = "concrete"
}

func (b *modernHouseBuilder) GetHouse() *house {
	b.buildRoof()
	b.buildWalls()
	b.buildFoundation()
	return &b.h
}

// конкретный объект
type oldHouseBuilder struct {
	h house
}

func (b *oldHouseBuilder) buildRoof() {
	b.h.roof = "tiles"
}

func (b *oldHouseBuilder) buildWalls() {
	b.h.walls = "wood"
}

func (b *oldHouseBuilder) buildFoundation() {
	b.h.foundation = "stone"
}

func (b *oldHouseBuilder) GetHouse() *house {
	b.buildRoof()
	b.buildWalls()
	b.buildFoundation()
	return &b.h
}

func main() {
	// каждый builder удовлетворяет houseBuilder, возвращает конкретную
	// имплементацию объекта
	var builders []houseBuilder
	builders = append(builders, &modernHouseBuilder{})
	builders = append(builders, &oldHouseBuilder{})

	for _, b := range builders {
		fmt.Println(b.GetHouse())
	}
}
