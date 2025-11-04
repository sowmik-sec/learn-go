package main

import "fmt"

func main() {
	var fruitList [4]string
	fruitList[0] = "Apple"
	fruitList[1] = "Banana"
	// fruitList[2] = "Tomato"
	fruitList[3] = "Peach"
	fmt.Println("Fruit list is: ", fruitList)

	var vegList = [3]string{"potato", "beans", "mushroom"}
	fmt.Println("Veg list is: ", vegList)
	var gadgetList = [4]string{"Phone", "laptop"}
	fmt.Println("Gadget list is: ", gadgetList)
}
