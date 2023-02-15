package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func binarySearch(nums []int, n int) (int, bool) {
	left := 0
	right := len(nums) - 1

	for left < right {
		mid := (left + right) / 2
		if nums[mid] < n {
			left = mid + 1
		} else if nums[mid] > n {
			right = mid - 1
		} else {
			return mid, true
		}

	}

	return 0, false
}

func main() {

	var nums [20]int
	n := 50

	rand.Seed(time.Now().Unix())
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	sort.Ints(nums[:])

	fmt.Printf("nums: %v\n", nums)
	if i, ok := binarySearch(nums[:], n); !ok {
		fmt.Printf("%v not found", n)
	} else {
		fmt.Printf("%v found at index %v", n, i)
	}
}
