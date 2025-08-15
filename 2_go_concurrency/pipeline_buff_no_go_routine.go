package main

import (
	"fmt"
)

func makeChannelBuff(numArr []int) <-chan int {

	// if the  6 is reduced dead lock will happen
	// if 6 it can take 7 values
	// always n + 1 vals
	mkChan := make(chan int, 6)

	// go func() {
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
	// }()

	return mkChan
}

func sqCovBuff(channel <-chan int) <-chan int {

	mkChan := make(chan int, 6)

	// go func() {
	for i := range channel {
		mkChan <- i * i
		fmt.Printf("Consuming Value --> %v \n", i)

		// fmt.Println(i * i)
	}
	close(mkChan)
	// }()

	return mkChan
}

func pipelineBuff() {

	// define an array
	num_arr := []int{1, 2, 3, 4, 5, 6}

	// Stage 1, create a channel
	numChannel := makeChannelBuff(num_arr)

	// Stage 2, sq the values
	sqArr := sqCovBuff(numChannel)

	// fmt.Println(sqArr)

	// // stage 3, convert to strings

	// strArr := stConv(sqArr)

	for i := range sqArr {
		fmt.Println(i)
	}
}
