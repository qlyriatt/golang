package main

import "fmt"

// каждая стратегия выполняет свое конкретное действие
type strategy interface {
	execute(int, int) int
}

type add struct {
}

func (st add) execute(a int, b int) int {
	return a + b
}

type sub struct {
}

func (st sub) execute(a int, b int) int {
	return a - b
}

type mult struct {
}

func (st mult) execute(a int, b int) int {
	return a * b
}

type context struct {
	strategy
}

func main() {

	var ctx context
	// в зависимости от выбранной стратегии,
	// контекст использует разные протоколы (совершает различные действия)

	ctx.strategy = add{}
	fmt.Println("add:", ctx.execute(5, 5))

	ctx.strategy = sub{}
	fmt.Println("sub:", ctx.execute(5, 5))

	ctx.strategy = mult{}
	fmt.Println("mult:", ctx.execute(5, 5))

}
