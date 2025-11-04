package main

import "fmt"

func main() {
	fmt.Println("Welcome to loops in golang")
	days := []string{"Sunday", "Tuesday", "Wednesday", "Friday", "Saturday"}
	fmt.Println(days)
	// for d := 0; d < len(days); d++ {
	// 	fmt.Println(days[d])
	// }
	// for i := range days {
	// 	fmt.Println(days[i])
	// }
	// for index, day := range days {
	// 	fmt.Printf("index is %v and value is %v\n", index, day)
	// }

	rougeValue := 0
	for rougeValue < 10 {
		rougeValue++
		if rougeValue == 5 {
			continue
		}
		if rougeValue == 7 {
			goto lco
		}
		fmt.Println("Value is: ", rougeValue)
	}
lco:
	fmt.Println("Here you are")
}
