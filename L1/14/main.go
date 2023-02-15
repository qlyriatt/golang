package main

import (
	"fmt"
	"reflect"
)

func printType(i any) {
	switch i.(type) {
	case int:
		fmt.Println("Int")
	case string:
		fmt.Println("String")
	case bool:
		fmt.Println("Bool")
	case chan int:
		fmt.Println("Chan Int")
	}
}

func main() {
	var i = []interface{}{40, "string", true, make(chan int)}

	// определение типа через switch type
	for _, val := range i {
		printType(val)
	}

	// определение типа через printf("%T")
	for _, val := range i {
		fmt.Printf("%T\n", val)
	}

	// определение типа через reflect.TypeOf()
	for _, val := range i {
		fmt.Print(reflect.TypeOf(val), "\n")
	}
}
