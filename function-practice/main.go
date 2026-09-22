package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("B must not be 0")
	}
	return a / b, nil
}

func main() {

	result1, err := divide(10, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result1)
	result2, err := divide(10, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result2)
}
