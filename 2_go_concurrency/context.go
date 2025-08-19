package main

import (
	"context"
	"sync"
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

}
