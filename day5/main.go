package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

// func maxOf(nums ...int) int {
// 	max := -11111111
// 	for _, val := range nums {
// 		if val > max {
// 			max = val
// 		}
// 	}
// 	if max == -11111111 {
// 		return 0
// 	}
// 	return max
// }

// func main() {
// 	nums := []int{4, 1, 9, 3}
// 	fmt.Println(maxOf(nums...))
// }

// func join(sep string, parts ...string) string {
// 	var sb strings.Builder
// 	for i, str := range parts {
// 		sb.WriteString(str)
// 		if i < len(parts)-1 {
// 			sb.WriteString(", ")
// 		}
// 	}
// 	return sb.String()
// }

// func main() {
// 	fmt.Println(join(", ", "a", "b", "c"))
// 	n, _ := strconv.Atoi("42")
// 	fmt.Println(n)
// }

func parseUser(line string) (name string, age int, err error) {
	person := strings.Split(line, ",")
	if len(person) < 2 {
		return "", 0, fmt.Errorf("invalid input format")
	}
	name = person[0]
	age, err = strconv.Atoi(strings.TrimSpace(person[1]))
	return
}

func main() {
	fmt.Println(parseUser("Alice,30"))
}
