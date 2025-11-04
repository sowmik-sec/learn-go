package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Welcome to files in golang")
	content := "This needs to go in a file"
	file, err := os.Create("./newFile.txt")
	checkNilErr(err)
	length, err := io.WriteString(file, content)
	checkNilErr(err)
	fmt.Println("Length is: ", length)
	defer file.Close()
	readFile("./newFile.txt")
}

func readFile(filename string) {
	databyte, err := ioutil.ReadFile(filename)
	checkNilErr(err)
	fmt.Println("Test data inside the file is \n", string(databyte))
}

func checkNilErr(err error) {
	if err != nil {
		panic(err)
	}
}
