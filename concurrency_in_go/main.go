package main

import "fmt"

func main() {
	var i int
	fmt.Scanln(&i)
	for x := 0; x < i; x++ {
		fmt.Println("Hello World")
	}
	//fmt.Println("Hello World")
}
