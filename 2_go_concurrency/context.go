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

	wg.Add(1)

	go func1(ctx, &wg, infiniteApples)

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
			// ok is the language provoded variable if stream ends its set to False
			case d, ok := <-stream:
				if !ok {
					fmt.Println("Channel CLosed")
					return
				}

				fmt.Println(d)
			}
		}
	}

	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)

	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}

	wg.Wait()

}
