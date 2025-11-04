package main

import "fmt"

func main() {
	fmt.Println("Maps in golang")
	languages := make(map[string]string)
	languages["JS"] = "Javascript"
	languages["RB"] = "Racing Bulls"
	languages["Py"] = "Python"
	fmt.Println("List of all languages: ", languages)
	fmt.Println("JS shorts for: ", languages["JS"])

	// loops
	for key, value := range languages {
		fmt.Printf("For key %v, value is %v\n", key, value)
	}
}
