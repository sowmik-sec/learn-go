// A1: get functions return type is int, error but it is returning only int. that's a problem. error must be returned.
// A2. we shouldn't use named return in large functions.
// A3. 0, false
// A4. (a) to statically ensure at compile time that a specific type implements a specific interface.
// (b) The statement import _ "image/png" has an effect because it triggers the package's init() function, which registers the PNG decoder into Go's central image registry.
// A5.
// a. package-level var initializers in main
// b. init() functions in dbdriver
// c. package-level var initializers in dbdriver
// d. main() body
// A6. it will panic with boom

package main

import "fmt"

func sum(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	result = a / b
	return
}

func main() {
	fmt.Println("Sum of discrete values: ", sum(1, 2, 3, 4))
	numbers := []int{10, 20, 30}
	fmt.Println("sum of spread slice: ", sum(numbers...))

	_, err := safeDivide(10, 0)
	if err != nil {
		fmt.Println("Error from safeDivide: ", err)
	} else {
		fmt.Println("no error")
	}
	double := func(n int) int {
		return n * 2
	}
	fmt.Println("Closure result (double 21): ", double(21))
}

func init() {
	fmt.Println("program started")
}
