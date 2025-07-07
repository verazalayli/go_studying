package concurrency

import (
	"fmt"
	"time"
)

/*
There are tasks (jobs) — they are put in the jobs channel.
There are workers — a fixed number of goroutines that:
   they read tasks from the jobs channel,
   processed,
   they send the result to the results channel (if necessary).
When all tasks are added, the jobs channel is closed.
After processing, all the workers are completed.
*/
// worker function — processes tasks from jobs and writes the result to results
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j)
		time.Sleep(time.Second) // симуляция работы
		fmt.Printf("Worker %d finished job %d\n", id, j)
		results <- j * 2
	}
}

func workerPool() {
	const numJobs = 5
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// run workers
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// sent jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // close channel — no more tasks

	// get results
	for a := 1; a <= numJobs; a++ {
		fmt.Println("Result:", <-results)
	}
}
