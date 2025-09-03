package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	name string `json:Name`
}

var userStore = make(map[int]User)

func main() {

	// multiplexer that allows us to control the flow of requets..

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", createUsers)

	// start the server

	fmt.Println("Server Started & Listening into port 8080.......")
	http.ListenAndServe(":8080", mux)
}

func handleRoot(
	w http.ResponseWriter,
	r *http.Request,

) {
	fmt.Fprintf(w, "Hello World")

}

func createUsers(
	w http.ResponseWriter,
	r *http.Request,
) {

	var newUser User

	// decode it to GO users
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(newUser)

}
