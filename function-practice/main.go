package main

import "fmt"

func sum(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func main() {
	fmt.Println(sum(1, 2, 3, 4))
	numbers := []int{1, 2, 3, 4}
	fmt.Println(sum(numbers...))
}
