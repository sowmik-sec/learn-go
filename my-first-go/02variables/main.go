package main

import "fmt"

const LoginToken string = "hqhchrweohhfsdf" // public variable

func main() {
	var username string = "Ahsan"
	fmt.Printf("Variable is of type: %T \n", username)

	var isLoggedIn bool = false
	fmt.Println(isLoggedIn)
	var smallVal uint8 = 23
	fmt.Println(smallVal)

	var anotherVariable int
	fmt.Println(anotherVariable)

	// implicit type
	var website = "example.com"
	fmt.Println(website)

	numberOfUsers := 30000
	fmt.Println(numberOfUsers)

	test := 3.1416
	fmt.Println(test)
}
