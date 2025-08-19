package main

// 0. Basic GO Routine

// func main() {
// 	// both gets executed concurrently.. if no go routine --> only hello gets printed
// 	go count("Hello", 0)
// 	count("Hafeel", 0)

// 	count("Hello", 0)
// 	count("Hafeel", 0)

// 	go count("Hello", 0)
// 	go count("Hafeel", 0)
// }

// 0.2 scenario without wait group
// func main() {

// 	go func() {
// 		// without wait group this func wont get executed.
// 		count("Sheep", 4000)
// 	}()

// 	// informs the GO Main Routine to wait until all the other GO Routines are stopped.
// }

// 0.3 Wait GRoups
// func main() {

// 	// define a waitGroup variable
// 	var wg sync.WaitGroup

// 	// Says need to wait for 1 Go Routine Variable
// 	wg.Add(1)

// 	go func() {
// 		count("Sheep", 4000)
// 		// says this GO Routine is done
// 		wg.Done()
// 	}()

// 	// informs the GO Main Routine to wait until all the other GO Routines are stopped.
// 	wg.Wait()
// }

// 0.4 The below code blocks the entire main thread causing a dead lock, because its a unbuffered pool so its SYNC, use buffered for ASYNC Perf

// func main() {
// 	c := make(chan string)
// 	c <- "hafeel"

// 	name := <-c
// 	fmt.Println(name)
// }

// 0.5 To resolve the above issue we need to use go routine

// func main() {
// 	c := make(chan string)
// 	var name string

// 	go func() {
// 		c <- "hafeel"
// 	}()

// 	// this blocks the main thread until it receives a message
// 	name = <-c

// 	fmt.Println(name)
// }

// 0.6 Buffered Pool ASYNC
// func main() {
// 	c := make(chan string, 3)
// 	c <- "hafeel"

// 	name := <-c
// 	fmt.Println(name)
// }

// 0.7 this code explains how the threads are executed
// here even thought the 500ms thread is ready we still need to wait 2 seconds for it to get its turn again
// func main() {
// 	c1 := make(chan string)
// 	c2 := make(chan string)

// 	go func() {
// 		for {
// 			c1 <- "Waiting for 500 milliseconds"
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			c2 <- "Waiting for 2 Seconds"
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()

// 	for {
// 		fmt.Println(<-c1)
// 		fmt.Println(<-c2)

// 	}
// }

// 0.8 The above scenario can be tackled by using SEELCT
// select prints the value from the channel that is ready to print..it wont wait liek it did in the prev 0.7 scenario

// func main() {
// 	c1 := make(chan string)
// 	c2 := make(chan string)

// 	go func() {
// 		for {
// 			c1 <- "Waiting for 500 milliseconds"
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			c2 <- "Waiting for 2 Seconds"
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()

// 	for {

// 		// whichever the channel is ready gets executed.
// 		select {
// 		case i := <-c1:
// 			fmt.Println(i)
// 			// time.Sleep(500 * time.Millisecond)

// 		case j := <-c2:
// 			fmt.Println(j)
// 			// time.Sleep(500 * time.Second)

// 		default:
// 			fmt.Println("Heelo")
// 		}
// 	}
// }

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

// ----------   undersatnding buffers  ------
// func main() {

// 	// 1.
// 	// unbuf()

// 	// 2.
// 	// infinite_routine()

// 	//3.
// 	// buff()

// 	// 4.
// 	// done_chnl()

// 	// // 5.
// 	// fmt.Println("---------------- Un Buffered Channel => START --------------")
// 	pipelineUnBuff()
// 	// fmt.Println("---------------- Un Buffered Channel => END --------------")

// 	// //6.
// 	// fmt.Println("---------------- Buffered Channel Without GO Routine => START --------------")
// 	// pipelineBuff()
// 	// fmt.Println("---------------- Buffered Channel Without GO Routine => END --------------")

// 	// // 7.
// 	// fmt.Println("---------------- Buffered Channel with GO Routine => START --------------")
// 	// pipelineBuffGoRoutine()
// 	// fmt.Println("---------------- Buffered Channel with GO Routine => END --------------")

// 	// 8. GO Routines Basics

// 	// chn_no_go_routine()
// 	// chn_go_routine()

// }

// -----------------------------------------------------------------------

// worker pools

// func main() {
// 	worker_pool_main()
// }

// ---------------------------

func main() {

	contextMain()
}
