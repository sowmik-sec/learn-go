package main

import "fmt"

type User struct {
	Name     string
	password string
}

func (u *User) CreateUser() {

}
func (u *User) validateUser() {

}

func main() {
	var user User
	user.CreateUser()
	a := new(int)
	b := make([]int, 5)
	c := make(map[string]int)

	fmt.Println(a) // I know these printing was not required.
	fmt.Println(b)
	fmt.Println(c)
}
