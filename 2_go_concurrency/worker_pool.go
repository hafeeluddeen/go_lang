package main

import "fmt"

func worker_pool_main() {

	jobs := make(chan int, 20)
	results := make(chan int, 20)

	// worker pools -> We can't gurantee the results will be in order but will be faster for sure.
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)

	// fill in the job
	for i := 0; i < 20; i++ {
		jobs <- i
	}

	// once we filled in the job channel we can close it
	close(jobs)

	// once we close the jobs we can print the results out

	for j := 0; j < 20; j++ {
		//If results is empty → the goroutine waits (blocks) until some other goroutine puts a value in.
		// this ensures all 20 values are pritned
		fmt.Println(<-results)
	}

}

// jobs channel is read only, results channel is for sending alone
// if we try to send on Jobs channel we will get compile time err, and read on results channel alsowe get error
func worker(jobs <-chan int, results chan<- int) {

	for i := range jobs {

		results <- fib(i)
	}

}

func fib(n int) int {

	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)

}
