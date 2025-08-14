package main

/* 1 GO Routine alone
func main() {

	// go routine fires this func execution to another tread
	go greeter("hello")
	go greeter("mate")
	go greeter("hope")

	// executed by the main thread.
	greeter("you are doing well")
}

func greeter(s string) {

	for i := 0; i < 3; i++ {

		time.Sleep(4 + time.Second)
		fmt.Println(s)
	}
}


*/

/* 2. GO Channels
func main() {

	// create a channel
	// channel is bascially a FIFO Q
	myChannel_1 := make(chan string)

	myChannel_2 := make(chan int)

	// Anon Func OR Nameless func
	go func() {

		// putting data into channel
		myChannel_1 <- "employee_1"

		myChannel_1 <- "employee_2"

	}()

	go func() {
		myChannel_2 <- 30
		myChannel_2 <- 40
	}()

	fmt.Println(myChannel_1) // op: 0xc00002a0e0 -> memory ref of the channel

	for i := 0; i < 2; i++ {

		// consume the channel
		// its a blocking line of code
		// the main go routine reads data from another go routine here.
		// the main func waits until it receives a message OR the channel is closed
		msg := <-myChannel_1
		msg_2 := <-myChannel_2

		fmt.Println(msg)
		fmt.Println(msg_2)

	}
}

*/

/*

3. SELECT


*/

// func main() {

// 	// create a channel
// 	// channel is bascially a FIFO Q
// 	myChannel_1 := make(chan string)

// 	myChannel_2 := make(chan int)

// 	// Anon Func OR Nameless func
// 	go func() {

// 		// putting data into channel
// 		myChannel_1 <- "employee_1"

// 		myChannel_1 <- "employee_2"

// 	}()

// 	go func() {
// 		myChannel_2 <- 30
// 		myChannel_2 <- 40
// 	}()

// 	/*
// 		Select blocks the code until one of the routine gets executed i.e meaning atleast one msg is returned.
// 		If multiple results are returned at same time.. GO chooses one at random.

// 	*/
// 	select {
// 	case msg_1 := <-myChannel_1:
// 		fmt.Println(msg_1)

// 	case msg_2 := <-myChannel_2:
// 		fmt.Println(msg_2)

// 	}

// }

/*

	Other stuff

*/

func main() {

	// 1.
	// unbuf()

	// 2.
	// infinite_routine()

	//3.
	buff()

	// 4.
	// done_chnl()

}
