package main

import "fmt"

func main() {
	fmt.Println("Welcome to functions in golang")
	greeter()

	greeterTwo()
	result := adder(3, 5)
	fmt.Println("Result is: ", result)
	proRes, proMsg := proAdder(3, 6, 2, 6)
	fmt.Println("Pro result is ", proRes)
	fmt.Println("Pro message is ", proMsg)
}

func proAdder(values ...int) (int, string) {
	total := 0
	for _, val := range values {
		total += val
	}
	return total, "Hello from pro adder"
}

func adder(valOne int, valTwo int) int {
	return valOne + valTwo
}

func greeter() {
	fmt.Println("Hello there")
}

func greeterTwo() {
	fmt.Println("another function")
}
