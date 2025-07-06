package concurrency

import (
	"fmt"
	"time"
)

/*
select is blocked, waiting for at least one channel to be ready (for reading or writing).
If several cases are ready, one of them is selected randomly.
If no channel is ready, but there is a default, it is executed immediately.
If there are no ready-made cases and default is missing, select blocks execution until ready.
*/
func selectFunc() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "from ch2"
	}()

	select {
	case msg1 := <-ch1:
		fmt.Println("Received:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received:", msg2)
	}
}
