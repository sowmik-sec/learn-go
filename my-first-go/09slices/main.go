package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Welcome to video on slices")
	var fruitList = []string{"potato", "beans", "mushroom"}
	fmt.Printf("Type of fruitList is %T\n", fruitList)
	fruitList = append(fruitList, "Mango", "Banana")
	fmt.Println(fruitList)
	fruitList = append(fruitList[1:3])
	fmt.Println(fruitList)

	highScores := make([]int, 4)
	highScores[0] = 234
	highScores[1] = 456
	highScores[2] = 678
	highScores[3] = 890
	fmt.Println(highScores)
	sort.Ints(highScores)
	fmt.Println(sort.IntsAreSorted(highScores))

	// how to remove a value from slices based on index
	var courses = []string{"reactjs", "javascript", "nextjs", "swift", "python", "ruby"}
	fmt.Println(courses)
	var index int = 2
	courses = append(courses[:index], courses[index+1:]...)
	fmt.Println(courses)
}
