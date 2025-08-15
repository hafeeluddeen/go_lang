package main

import (
	"fmt"
	"time"
)

func makeChannel(numArr []int) <-chan int {

	mkChan := make(chan int)

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
			time.Sleep(5 * time.Millisecond)

		}

		close(mkChan)
	}()

	return mkChan
}

func sqCov(channel <-chan int) <-chan int {

	mkChan := make(chan int)

	go func() {
		for i := range channel {
			mkChan <- i * i
			fmt.Printf("Consuming Value --> %v \n", i*i)
			time.Sleep(5 * time.Millisecond)

		}
		close(mkChan)
	}()

	// for i := range channel {
	// 	mkChan <-  i* i
	// 	fmt.Printf("Consuming Value --> %v \n", i*i)
	// 	time.Sleep(5 * time.Millisecond)

	// }

	return mkChan
}

// func stConv(numArr <-chan string) []string {

// 	stArr := []string{}
// 	for _, i := range numArr {
// 		stArr = append(stArr, strconv.Itoa(i))
// 	}

// 	return stArr
// }

func pipelineUnBuff() {

	// define an array
	num_arr := []int{1, 2, 3, 4, 5, 6}

	// Stage 1, create a channel
	numChannel := makeChannel(num_arr)

	// Stage 2, sq the values
	sqArr := sqCov(numChannel)

	// fmt.Println(sqArr)

	// // stage 3, convert to strings

	// strArr := stConv(sqArr)

	for i := range sqArr {
		fmt.Printf("Printing Value --> %v \n", i)
		time.Sleep(5 * time.Millisecond)

	}
}
