package main

import "fmt"

func add(lhs, rhs int) int {
	return lhs + rhs
}


func compute(lhs, rhs int, op func(lhs, rhs int) int) int {
	fmt.Printf("Running a computation with %v & %v\n", lhs, rhs)
	return op(lhs, rhs)
}
