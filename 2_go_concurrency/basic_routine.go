package main

import "fmt"

func count(name string, counter int) {

	if counter == 0 {
		counter = 100000000
	}

	for i := 0; i < counter; i++ {
		fmt.Println(name)
	}
}
