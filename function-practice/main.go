package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func main() {
	var operation func(int, int) int
	operation = add
	fmt.Println(operation(10, 20))
	multiply := func(a, b int) int {
		return a * b
	}
	fmt.Println(multiply(10, 20))
}
