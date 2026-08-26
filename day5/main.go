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

// func clamp(n, lo, hi int) (result int) {
// 	if n < lo {
// 		result = lo
// 		return
// 	} else if n > hi {
// 		result = hi
// 		return
// 	}
// 	return
// }

// func main() {
// 	fmt.Println(clamp(15, 0, 10))
// 	fmt.Println(clamp(-3, 0, 10))
// }

func maxOf(nums ...int) int {
	max := -11111111
	for _, val := range nums {
		if val > max {
			max = val
		}
	}
	if max == -11111111 {
		return 0
	}
	return max
}

func main() {
	nums := []int{4, 1, 9, 3}
	fmt.Println(maxOf(nums...))
}
