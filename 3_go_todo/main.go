package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
)

type Task struct {
	task_no      string    `csv:"task_no"`
	task_title   string    `csv:"task_title"`
	created_date time.Time `csv:"created_date"`
	due_date     string    `csv:"due_date"`
	is_completed bool      `csv:"is_completed"`
}

func main() {

	// run for ever
	for {

		var option int

		fmt.Println("Choose any one of the option")
		fmt.Println("----------------------------------------")
		fmt.Println("0. Display Tasks")
		fmt.Println("1. Add Tasks")
		fmt.Println("2. Mark a task as complete")
		fmt.Println("3. Delte Task")
		fmt.Println("----------------------------------------")
		fmt.Scanln(&option)

		// Add a task
		if option == 0 {
			display_tasks()
		} else if option == 1 {
			add_task()

		} else if option == 2 {
			mark_as_done()
		} else {
			fmt.Println("Invalid Option Choose Prop")
		}

	}

}

func create_csv() {

	// define the headers
	header := []string{"task_no", "task_title", "created_date", "due_date", "is_completed"}

	// create a csv
	file, err := os.Create("tasks.csv")

	if err != nil {
		log.Fatalf("Failed to create csv!!")
	}

	// close the file when you exit this function
	defer file.Close()

	// create a csv writer
	writer := csv.NewWriter(file)

	// write at the end of file
	defer writer.Flush()

	// write into the csv
	if err := writer.Write(header); err != nil {
		log.Fatalln("Failed to write header")
		return
	}

	fmt.Println("Created File successfully")
}

// check if the file exsists if not just creates the file
func fileExists() {

	_, err := os.Stat("tasks.csv")

	if err != nil {
		// log.Fatalf("Failed to check if file exists")
		// return false
		create_csv()
	}

}

func display_tasks() {

	fileExists()
	/*
		os.ModePerm → shorthand for 0777 (read, write, execute permissions for everyone).

		This only matters if the file is created. If it already exists, the existing permissions are kept.
	*/
	file, err := os.OpenFile("tasks.csv", os.O_CREATE|os.O_RDWR, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to open the file")
	}
	/*
		When we do CLOSE the file the lock gets released.
		If we don't close the file the LOCK will not be released..It will lock the file forever
		No matter how the file exits defer will always be executed
		if ther are 3 deferes in a func..Their order of execution id backwards
		defer_3 -> defer_2 -> defer_1
	*/
	defer file.Close()

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		log.Fatalf("File in use")
		return
	}

	// create the csv reader
	reader := csv.NewReader(file)

	// read all row into the memory
	data, err := reader.ReadAll()

	if err != nil {
		log.Fatalf("Failed to read the file")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header row
	fmt.Fprintf(w, "TaskNo\tTitle\tCreatedDate\tDueDate\tCompleted\tTimeDiff\n")

	for _, row := range data[1:] {
		// parse the created date
		t, err := time.Parse(time.RFC3339, row[2])
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue // skip this row instead of exiting
		}

		// compute how far behind it is
		diff := timediff.TimeDiff(t)

		// write to tabwriter
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
			row[0], row[1], row[2], row[3], row[4], diff)
	}

	// Flush at the end
	w.Flush()

	// time.Sleep(15 * time.Second)

}

func add_task() {
	// open the file in write mode
	file, err := os.OpenFile("tasks.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		log.Fatalf("Failed to open the file")
	}

	defer file.Close()

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return
	}

	// create a csv writer
	writer := csv.NewWriter(file)

	// write at the end of file
	defer writer.Flush()

	// get ip

	var task_name string
	fmt.Println("Enter Task Name")
	fmt.Scanln(&task_name)

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalln("Failed to get the time zone")
		return
	}
	now := time.Now().In(ist)

	var due_date string
	fmt.Println("Enter due date")
	fmt.Scanln(&due_date)

	is_completed := false

	// construct the task
	t := Task{
		task_no:      "99",
		task_title:   task_name,
		created_date: now,
		due_date:     due_date,
		is_completed: is_completed,
	}

	if err := writer.Write([]string{
		t.task_no,
		t.task_title,
		t.created_date.Format(time.RFC3339),
		t.due_date,
		strconv.FormatBool(t.is_completed)}); err != nil {
		log.Fatalln("Faield to write")
	}

	// time.Sleep(15 * time.Second)

}

func mark_as_done() {
	var option string

	fmt.Println("Which taks u had completed hurray!!")
	fmt.Scanln(&option)

	// open and read the file

	file, err := os.Open("tasks.csv")

	if err != nil {
		log.Fatalln("Failed to open the fle")
		return
	}

	defer file.Close()

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		return
	}

	reader := csv.NewReader(file)

	// bring the whole daat to memory to edit it.
	rows, err_2 := reader.ReadAll()

	if err_2 != nil {
		log.Fatalf("Failed to open the file err_2")
	}

	// iterate through all the rows and update
	for i, row := range rows {
		if row[0] != option {
			continue
		}

		rows[i][4] = "true"
	}

	// finally wrtite into the csv
	// jst overwritting the data not efficient actually

	file, err = os.Create("tasks.csv")

	if err != nil {
		log.Fatalln("Failed to open the fle")
		return
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(rows); err != nil {
		log.Fatalln("Failed to update the fle")
		return
	}

	// open to edit the file
}

/*

	1. mutex lock when writting to a file  --> Done
	2. implement timestamp properly    --> done
	3. Add cli
	4. Add Tidy my desk   --> done

*/
