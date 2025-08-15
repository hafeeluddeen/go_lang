package main

import (
	"fmt"
)

func makeChannelBuffGoRoutine(numArr []int) <-chan int {

	mkChan := make(chan int, 3)

	go func() {
		/*
			1.
			since its an unbuff channel
			once a single value is put into it, it goes into a blocking mode, cauz unbuff is SYNC in nature

			2.
			only when the channel in sqConv consumes a value, this will put the next value into the channel until then its blocked

		*/
		for _, i := range numArr {
			mkChan <- i
			fmt.Printf("Putting Value --> %v \n", i)
		}

		close(mkChan)
	}()

	return mkChan
}

func sqCovBuffGoRoutine(channel <-chan int) <-chan int {

	mkChan := make(chan int, 6)

	go func() {
		for i := range channel {
			mkChan <- i * i
			fmt.Printf("Consuming Value --> %v \n", i)
		}
		close(mkChan)
	}()

	return mkChan
}

func pipelineBuffGoRoutine() {

	// define an array
	num_arr := []int{1, 2, 3, 4, 5, 6, 7}

	// Stage 1, create a channel
	numChannel := makeChannelBuffGoRoutine(num_arr)

	// Stage 2, sq the values
	sqArr := sqCovBuffGoRoutine(numChannel)

	// fmt.Println(sqArr)

	// // stage 3, convert to strings

	// strArr := stConv(sqArr)

	for i := range sqArr {
		fmt.Println(i)
	}
}
