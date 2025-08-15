package main

import (
	"fmt"
	"time"
)

// 1.
func buff() {

	// define the Buffered channel
	myChannel := make(chan int, 3)

	// define int array
	int_arr := []int{21, 12, 2002}

	for _, i := range int_arr {

		select {

		// put value into the channel
		case myChannel <- i:

		}
	}

	// close the Channel -> No more accepts value
	close(myChannel)

	// consume the channel
	// GO will automatically cut itself after reading 3 values because it knows where the channel ends
	for i := range myChannel {
		fmt.Println(i)
	}

}

// 2.

func infinite_routine() {

	// Select stmt combined with Anon func helps the GO Routine to RUN forever.
	go func() {
		for {
			select {
			default:
				fmt.Println("Doing Work")
			}
		}
	}()

	time.Sleep(5 * time.Second)
}

// 3.

func unbuff() {

	// since unbuffered channels are SYNC in nature, its like their size is 1
	ch := make(chan int)

	go func() {
		ch <- 1 // blocks until buff goroutine reads
		fmt.Println("sent 1")
		ch <- 2
		fmt.Println("sent 2")
	}()

	time.Sleep(time.Second) // sender is stuck here until we read
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)

	time.Sleep(5 * time.Second)

}

// 4.

// DONE channel

// this says the func accepts a chan of type bool
// Also this syntax says the channel is read only
func doWork(done <-chan bool) {

	for {

		// the "SELECT" here is extremely useful as it helps us to check if the parent wants to cancel the channel
		// once command comes through we can cancel it
		select {

		// the moment we consume something from our channel we will wrap up our work
		// once we close the channel select ackonwledges it and exits.
		// we are consuming done here
		case <-done:
			fmt.Print(<-done)
			return

		default:
			fmt.Println("Doing Work")

		}
	}
}

func done_chnl() {

	// need to run a routine for some desried time and stop it when we wish to stop it.

	// define a channel
	done := make(chan bool)

	go doWork(done)

	time.Sleep(5 * time.Second)

	// close the channel
	close(done)

}
