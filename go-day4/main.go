package main

import "fmt"

func increase(p *int) {
	*p++
}

func main() {
	var p *int
	x := 10
	p = &x
	increase(p)
	fmt.Println(*p)
}
