package main

import "fmt"

func main() {
	fmt.Println("Structs in golang")
	sowmik := User{"Sowmik", "ex@example.com", true, 53}
	fmt.Println(sowmik)
	fmt.Printf("Sowmik details are: %v\n", sowmik)
	fmt.Printf("Name is %v, email is %v\n", sowmik.Name, sowmik.Email)
}

type User struct {
	Name   string
	Email  string
	Status bool
	Age    int
}
