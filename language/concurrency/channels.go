package concurrency

import "fmt"

/*
Channels allow one goroutine to send data and the other to receive it,
ensuring secure transmission of information without the need for manual synchronization.
*/
func channel() {
	ch := make(chan string) // channel for string
	//ch <- 42      // отправка значения 42 в канал
	//val := <-ch   // получение значения из канала

	go func() {
		ch <- "hello from goroutine" // отправка в канал
	}()

	msg := <-ch // блокируемся до получения
	fmt.Println(msg)
	close(ch) // Closing channel
	// After closing we can only read from channel
	for val := range ch {
		fmt.Println(val)
	}
}

// Unbuffered channels used by default
// Buffered channels:
func unbufChannels() {
	ch := make(chan int, 3)
	ch <- 1 // не блокирует, если буфер не полон
}

// Send-only / receive-only channels
func sendOnly(ch chan<- int) {
	ch <- 10
}

func receiveOnly(ch <-chan int) {
	fmt.Println(<-ch)
}

/*
| Problem                    | Description                                                                     |
| -------------------------- | ------------------------------------------------------------------------------- |
| 🛑 **Deadlock**            | When all the goroutines are waiting and no one is doing the sending or reading. |
| 🧵 **Excess goroutines**   | Don't forget to complete the goroutines and close the channels.                 |
| ❌ **Reading from closed** | Acceptable, but returns a 0-type value, `ok == false'.                          |
*/
