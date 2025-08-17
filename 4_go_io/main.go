package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	writer_main()
}

// Reader
func str_reader() {

	/*

		The function count() can be reused here once again because string.NewReader also implemnts io.Reader, so no code chnages are required.


	*/
	r := strings.NewReader("Hello World")
	res, err := count(r)

	if err != nil {
		log.Fatalf("Failed to return file counter")
	}

	fmt.Printf("Total no of words -> %v \n", res)

}

func read_frm_file() {

	file, err := os.Open("test.txt")
	defer file.Close()
	/*
		os.Open => return *os.File which internally implements io.Reader


	*/

	if err != nil {
		log.Fatalf("Failed to open the file")
	}

	res, err := count(file)

	if err != nil {
		log.Fatalf("Failed to return file counter")
	}

	fmt.Printf("Total no of words -> %v \n", res)

}

func count(ptr io.Reader) (int, error) {

	count := 0
	buffer := make([]byte, 1024)

	/*

		Need to use for loop because it reads only first n bytes from memory..if there are 2000 bytes and no for loop result will br 2000

		Try commenting for loop and reducing buffer size than actual file size

	*/
	for {

		// puts n elements in to byte array buffer
		n, err := ptr.Read(buffer)

		for _, l := range buffer[:n] {
			if (l >= 'A' && l <= 'Z') || (l >= 'a' && l <= 'z') {
				fmt.Printf("%c", l)
				count++
			}

			// if l >= '1' && l <= '9' {
			// 	count++
			// }

		}

		if err == io.EOF {
			return count, nil
		}

		if err != nil {
			return 0, err
		}
	}

	// return count, nil
}

// writter

func writer_main() {
	/*

		os.Create implements io.Wrtiter at the core.

	*/
	file, err := os.Create("writter.txt")

	if err != nil {
		log.Fatalf("Failed to create file")
	}

	defer file.Close()

	// write contents into the file

	n, err := file_writter("Hello World !!", file)

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("File written Successfullt no of bytes: %v \n", n)
}

// func file_writter(string s, w io.Writer) (int, error){
// 	// create a byte array

// 	n, err := w.Write([]byte(s))

// 	if err != nil{
// 		return 0,log.Fatalf("Failed to wrtie into the file")
// 	}

// }

func file_writter(s string, w io.Writer) (int, error) {
	// create a byte array

	n, err := w.Write([]byte(s))

	if err != nil {
		// 		return 0,log.Fatalf("Failed to wrtie into the file")

		return 0, fmt.Errorf("Failed to wrtie into the file")
	}

	return n, nil

}
