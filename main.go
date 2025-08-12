package main

import (
	"fmt"
	"time"
)

func main() {
	x := make(map[int]int, 1)
	go func() { x[1] = 2 }()
	go func() { x[3] = 7 }()
	go func() { x[123] = 10 }()
	go func() { x[1] = 2 }()
	go func() { x[34] = 7 }()
	go func() { x[1432] = 10 }()
	go func() { x[1] = 2 }()
	go func() { x[100] = 7 }()
	go func() { x[34] = 10 }()
	go func() { x[1] = 2 }()
	time.Sleep(10 * time.Second) //блокируемся на 100 миллисекунд
	fmt.Println("x[1] =", x[1])
}
