package main

import "fmt"

func main() {
	x := 10
	p := &x
	fmt.Println(x)
	fmt.Println(p)
	fmt.Println(*p)
	x = 20
	fmt.Println(*p)
	*p = 30
	fmt.Println(x)
}
