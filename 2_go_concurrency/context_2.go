package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
	CONTEXT is mom and child proocess are like actual children..they will listen to mom, if mom says done they all will come home

*/

func ctxMain() {

	// i need something that can control my 3 child porcess
	// for that am going with context
	ctx, cancel := context.WithCancel(context.Background())

	// at some poit when the func ends need to cancel the context so use defer
	defer cancel()

	// since i have 3 child process with me i need to wait for them to finish.. so am using waitGroups.

	/*
		waitGroups are used to wait for the child to finish..
		context is like a remote control instructuing those child process like what to do ( stop, start) so on.
	*/
	var wg sync.WaitGroup

	// 1. need a func for generating values continously and put values in to the buffer.

	// its job is just to put values into the buffer non stop
	generator := func(data string, buf chan interface{}) {

		for {

			select {

			// check if the context says finsish all work and come home
			/*
				CONTEXT is mom and child proocess are like actual children..they will listen to mom, if mom says done they all will come home
			*/

			case <-ctx.Done():
				return

			// if not done just put value into the buffer
			case buf <- data:

			}
		}

	}

	// in this examle am assuming my parent (main) will have 3 child
	// so creating 3 buff where values are putinto it continously

	// Generator for Child 1
	// this buffer will store infin apples into it.
	infinStarberry := make(chan interface{})
	go generator("Starberry", infinStarberry)

	// Generator for child 2
	infinApples := make(chan interface{})
	go generator("Apples", infinApples)

	// Generator for child 3
	infinLemon := make(chan interface{})
	go generator("Lemon", infinLemon)

	wg.Add(1)
	child1(ctx, &wg, infinStarberry)

	// wg.Add(1)
	// go otherChilds(ctx, &wg, infinApples)

	// wg.Add(1)
	// go otherChilds(ctx, &wg, infinLemon)

	// since this child process has its own child, need to say to this wait for them
	wg.Wait()

}

// child 1 consume
/*
	parentCtx -> Its the parent remote control given to the child

	parentWg -> Parents wait group, once this child has done with its thing it can say to the parent its done.

	stream -> buffer from where this guy needs to consume
*/
func child1(parentCtx context.Context, parentWg *sync.WaitGroup, stream <-chan interface{}) {

	// this child when finishig its task must inform its parent that its done with its work and the parent need not wait for this now
	// Since this func itself a child Process of main process we are using wait_group.Done() to tell its parent once its done.
	defer parentWg.Done()

	/*
		As per the Architecture of this example the child have 3 more childs itself.

		Since the Child is having child of its own to monitor its child we need a wait group for the child

				Parent
				  |
				  v
				 child 1
				  |
				  V
			  child_child_1
	*/

	var wg sync.WaitGroup

	// need to create one remote control for controlling our childern
	// So using context here
	// To control its child (child_child_1), the child (child_1) needs an context
	// ctx.Timeout -> takes another context as a parameter and the time after whichit should elapse
	newCtx, cancel := context.WithTimeout(parentCtx, time.Second*5)
	// at somepoint this context must end so using defer
	defer cancel()

	// func for creating childs child Process

	childsChildProcess := func(ctx context.Context) {

		// this child when finishig its task must inform its parent that its done with its work and the parent need not wait for this now
		// since we know its a child process we are using waait_group.Done() to let it infomr its parent
		defer wg.Done()

		// now this child's child process needs to consume the value from the infinite buffer
		for {
			select {
			case <-ctx.Done():
				return

			case d, ok := <-stream:
				if !ok {
					fmt.Println("Stream Ended!!")
					return
				}

				fmt.Println(d)

			}
		}
	}

	// Start the child's child Process
	for i := 0; i < 3; i++ {
		wg.Add(1)
		// subroutine -> child process
		go childsChildProcess(newCtx)

	}

	// since this child process has its own child, need to say to this wait for them
	wg.Wait()

}

func otherChilds(ctx context.Context, parentWg *sync.WaitGroup, stream <-chan interface{}) {

	defer parentWg.Done()

	// consume the stream infinitely
	for {

		select {
		case <-ctx.Done():
			return
		case d, ok := <-stream:
			if !ok {
				fmt.Println("Stream Closed")
				return
			}

			fmt.Println(d)
			// time.Sleep(time.Millisecond * 100)
		}

	}

}
