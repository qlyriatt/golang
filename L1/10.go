package main

import (
	"fmt"
	"math"
)

type tempGroup struct {
	tempRange int
	temps     []float64
}

// проверка принадлежности temp к шагу tempRange
func doesInclude(positive bool, tempRange int, temp float64) bool {
	if positive {
		t := float64(tempRange+10) - temp
		return 10 > t && t > 0
	} else {
		t := float64(tempRange-10) - temp
		return -10 < t && t < 0
	}
}

// проверка наличия группы с шагом tempRange
func tempRangeExists(tempRange int) bool {
	for _, val := range tempGroups {
		if val.tempRange == tempRange {
			return true
		}
	}
	return false
}

var temps = [...]float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
var tempGroups []tempGroup

func main() {

	for _, temp := range temps {
		var t tempGroup
		if temp >= 0 { // для положительных температур группы 0..10 , 10..20
			t.tempRange = int(math.Floor(temp/10)) * 10
			if tempRangeExists(t.tempRange) {
				continue
			}
			for _, temp := range temps {
				if doesInclude(true, t.tempRange, temp) {
					t.temps = append(t.temps, temp)
				}
			}

		} else { // для отрицательных температур группы 0..-10 , -10..-20
			t.tempRange = int(math.Ceil(temp/10)) * 10
			if tempRangeExists(t.tempRange) {
				continue
			}
			for _, temp := range temps {
				if doesInclude(false, t.tempRange, temp) {
					t.temps = append(t.temps, temp)
				}
			}
		}
		tempGroups = append(tempGroups, t)
	}

	for _, val := range tempGroups {
		fmt.Printf("%v:%v\n", val.tempRange, val.temps)
	}
}
