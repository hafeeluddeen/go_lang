package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func contextMain() {

	var wg sync.WaitGroup

	/*
		Bascially cancel func closes the done channel

	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := func(dataItem string, stream chan interface{}) {
		for {
			select {
			case <-ctx.Done():
				return

			case stream <- dataItem:
			}
		}
	}

	infiniteApples := make(chan interface{})
	go generator("Apple", infiniteApples)

	infiniteBananas := make(chan interface{})
	go generator("Bananas", infiniteBananas)

	infiniteCheerys := make(chan interface{})
	go generator("Cherrys", infiniteCheerys)

	func2 := genericFunc
	func3 := genericFunc

	wg.Add(1)
	go func1(ctx, &wg, infiniteApples)

	wg.Add(1)
	go func2(ctx, &wg, infiniteBananas)

	wg.Add(1)
	go func3(ctx, &wg, infiniteCheerys)

	wg.Wait()

}

func func1(ctx context.Context, parentWg *sync.WaitGroup, stream <-chan interface{}) {

	defer parentWg.Done()

	var wg sync.WaitGroup

	doWork := func(ctx context.Context) {

		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			// ok is the language provided variable if stream ends its set to False
			case d, ok := <-stream:
				if !ok {
					fmt.Println("Channel CLosed")
					return
				}

				fmt.Println(d)
			}
		}
	}

	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)

	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}

	wg.Wait()

}

func genericFunc(ctx context.Context, wg *sync.WaitGroup, stream <-chan interface{}) {
	defer wg.Done()

	// it keeps on running until it receives a val from the channel
	for {
		select {
		case <-ctx.Done():
			return

		case d, ok := <-stream:
			if !ok {
				fmt.Println("Channel CLosed")
				return
			}
			fmt.Println(d)
		}

	}
}
