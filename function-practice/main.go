package main

import "fmt"

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func main() {
	fmt.Println(minMax(4, 10))
}
