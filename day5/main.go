package main

import "fmt"

// func half(n int) (result int, err bool) {
// 	if n%2 != 0 {
// 		err = false
// 		return
// 	} else {
// 		result = n / 2
// 		err = true
// 		return
// 	}
// }

// func main() {
// 	fmt.Println(half(3))
// }

func clamp(n, lo, hi int) (result int) {
	if n < lo {
		result = lo
		return
	} else if n > hi {
		result = hi
		return
	}
	return
}

func main() {
	fmt.Println(clamp(15, 0, 10))
	fmt.Println(clamp(-3, 0, 10))
}
