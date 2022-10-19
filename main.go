package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/post-project", home).Methods("GET")
	r.HandleFunc("/add-project", addProject).Methods("GET")
	r.HandleFunc("/detail-project", detailProject).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/post-project", postMyProject).Methods("POST")

	fmt.Println("Server is running on port 5678...\t(press \"ctrl + c\" to cancel)")
	http.ListenAndServe("localhost:5678", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/my-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/detail-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func postMyProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Project Name = " + r.PostForm.Get("input-project-name"))
	fmt.Println("Start Date = " + r.PostForm.Get("input-start-date"))
	fmt.Println("End Date = " + r.PostForm.Get("input-end-date"))
	fmt.Println("Email = " + r.PostForm.Get("input-email"))
	fmt.Println("React Js = " + r.PostForm.Get("check-reactjs"))
	fmt.Println("Vue Js = " + r.PostForm.Get("check-vuejs"))
	fmt.Println("Angular = " + r.PostForm.Get("check-angular"))
	fmt.Println("Laravel = " + r.PostForm.Get("check-laravel"))

	http.Redirect(w, r, "/post-project", http.StatusMovedPermanently)
}
