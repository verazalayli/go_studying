package main

import "fmt"

func main() {
	first := []int{10, 20, 30, 40}
	second := make([]*int, len(first))
	for i, v := range first {
		second[i] = &v
	}
	fmt.Println(second)
	fmt.Println(*second[0], *second[1])
}
