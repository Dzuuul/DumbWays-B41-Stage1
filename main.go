package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/add-project", addProject).Methods("GET")
	r.HandleFunc("/detail-project", detailProject).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/post-project", postMyProject).Methods("POST")

	fmt.Println("Server port 5678 activated")
	http.ListenAndServe("localhost:5678", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/my-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/detail-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}

func postMyProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}
