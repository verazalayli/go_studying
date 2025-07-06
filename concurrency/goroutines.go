package concurrency

import (
	"fmt"
	"time"
)

/*
	go use 3-level model:

G — Goroutine
M — ОС stream (Machine)
P — Processor (logic processor Go, limit quantity
*/
func sayHello() {
	fmt.Println("Hello from goroutine")
}

func goroutines() {
	go sayHello()               // run in the background
	time.Sleep(1 * time.Second) // give time to run goroutine
}

// Using channels
func worker1(ch chan string) {
	ch <- "done" // sent
}

func goroutinesWithChannel() {
	ch := make(chan string)
	go worker1(ch)
	msg := <-ch // get
	fmt.Println(msg)
}

// Leak goroutines
func leak() {
	ch := make(chan string)
	go func() {
		fmt.Println(<-ch) // block forever
	}()
}

/*
	Goroutines pulls

Goroutine pools are a pattern that limits the number of simultaneously running goroutines.
Why do I need a goroutine pool?
Goroutines are easy, but if you run hundreds of thousands at the same time:
A lot of memory is consumed,
Load on the Go Scheduler,
Leaks and deadlocks are possible,
The CPU is heavily loaded.
*/
func worker2(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		results <- j * 2
	}
}

func goroutinesPools() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for w := 1; w <= 3; w++ {
		go worker2(w, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 5; a++ {
		fmt.Println(<-results)
	}
}
