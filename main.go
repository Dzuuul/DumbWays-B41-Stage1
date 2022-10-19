package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloWorld)
	r.HandleFunc("/home", helloHome)
	r.HandleFunc("/project", helloProject)

	fmt.Println("Server port 5678 activated")
	http.ListenAndServe("localhost:5678", r)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World! Welcome back to my channel!"))
}

func helloHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Home! Nice to meet you!"))
}

func helloProject(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Project! I think you need additional projects!"))
}
