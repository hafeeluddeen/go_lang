package main

import "fmt"

func chn_no_go_routine() {
	ch := make(chan int) // unbuffered

	// blocks the entire main thread
	// results in deadlock
	ch <- 1 // send

	// this will be never reached cauz the program stucks in sync mode
	fmt.Println("Sent 1")

}

func chn_go_routine() {
	ch := make(chan int)

	// starts a seperate sub-thread
	go func() {
		ch <- 1
		fmt.Println("Sent 1")
	}()

	/*
		if the receiver is commented out "Sent 1" will not be printed
		because the main thread complets before the GO Routine gets executed.
	*/

	// receiver
	// fmt.Println(<-ch)

}
