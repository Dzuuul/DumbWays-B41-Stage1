package main

import (
	"context"
	"day-9/connection"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Project struct {
	Id              int
	Project_name    string
	Start_date      time.Time
	End_date        time.Time
	Detail_duration string
	Duration        string
	Description     string
	Technologies    []string
}

var dataProject = []Project{}

func main() {
	r := mux.NewRouter()

	connection.DatabaseConnect()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/home", project).Methods("GET")
	r.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	r.HandleFunc("/form-edit-project/{index}", formEditProject).Methods("GET")
	r.HandleFunc("/detail-project/{index}", detailProject).Methods("GET")
	r.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	r.HandleFunc("/add-my-project", addProject).Methods("POST")
	r.HandleFunc("/edit-my-project/{index}", editProject).Methods("POST")

	fmt.Println("Server is running on port 5656...\t(press \"ctrl + c\" to abort)")
	http.ListenAndServe("localhost:5656", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/home.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, duration, description, technologies FROM tb_projects ORDER BY id DESC")

	var result []Project
	for data.Next() {
		var each = Project{}

		err = data.Scan(&each.Id, &each.Project_name, &each.Duration, &each.Description, &each.Technologies)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	respData := map[string]interface{}{
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	ProjectDetail := Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				Project_name:    data.Project_name,
				Detail_duration: data.Detail_duration,
				Duration:        data.Duration,
				Description:     data.Description,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-me.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	inputProjectName := r.PostForm.Get("input-project-name")
	inputStartDate := r.PostForm.Get("input-start-date")
	inputEndDate := r.PostForm.Get("input-end-date")
	inputDescription := r.PostForm.Get("input-description")

	var inputTechnologies []string
	inputTechnologies = r.Form["technologies"]
	fmt.Println(inputTechnologies)

	parseStartDate, _ := time.Parse("2006-01-02", inputStartDate)
	parseEndDate, _ := time.Parse("2006-01-02", inputEndDate)

	hour := parseEndDate.Sub(parseStartDate).Hours()
	day := hour / 24
	week := day / 7
	month := day / 30
	year := day / 365

	formatStartDate := parseStartDate.Format("2 Jan 2006")
	formatEndDate := parseEndDate.Format("2 Jan 2006")

	inputDetailDuration := formatStartDate + " - " + formatEndDate

	var inputDuration string

	switch {
	case year == 1:
		inputDuration = strconv.FormatFloat(year, 'f', 0, 64) + " year"
	case year > 1:
		inputDuration = strconv.FormatFloat(year, 'f', 0, 64) + " years"
	case month == 1:
		inputDuration = strconv.FormatFloat(month, 'f', 0, 64) + " month"
	case month > 1:
		inputDuration = strconv.FormatFloat(month, 'f', 0, 64) + " months"
	case week == 1:
		inputDuration = strconv.FormatFloat(week, 'f', 0, 64) + " week"
	case week > 1:
		inputDuration = strconv.FormatFloat(week, 'f', 0, 64) + " weeks"
	case day == 1:
		inputDuration = strconv.FormatFloat(day, 'f', 0, 64) + " day"
	case day > 1:
		inputDuration = strconv.FormatFloat(day, 'f', 0, 64) + " days"
	default:
		inputDuration = "WRONG DATE!"
	}

	newProject := Project{
		Project_name:    inputProjectName,
		Start_date:      parseStartDate,
		End_date:        parseEndDate,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
		Technologies:    inputTechnologies,
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-my-project.html")

	if err != nil {
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/edit-my-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	edit := Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			edit = Project{
				Project_name:    data.Project_name,
				Start_date:      data.Start_date,
				End_date:        data.End_date,
				Detail_duration: data.Detail_duration,
				Duration:        data.Duration,
				Description:     data.Description,
			}
		}
	}

	data := map[string]interface{}{
		"Index": index,
		"Edit":  edit,
	}

	tmpl.Execute(w, data)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	inputProjectName := r.PostForm.Get("input-project-name")
	inputStartDate := r.PostForm.Get("input-start-date")
	inputEndDate := r.PostForm.Get("input-end-date")
	inputDescription := r.PostForm.Get("input-description")

	var inputTechnologies []string
	inputTechnologies = r.Form["technologies"]
	fmt.Println(inputTechnologies)

	parseStartDate, _ := time.Parse("2006-01-02", inputStartDate)
	parseEndDate, _ := time.Parse("2006-01-02", inputEndDate)

	hour := parseEndDate.Sub(parseStartDate).Hours()
	day := hour / 24
	week := math.Round(day / 7)
	month := math.Round(day / 30)
	year := math.Round(day / 365)

	formatStartDate := parseStartDate.Format("2 Jan 2006")
	formatEndDate := parseEndDate.Format("2 Jan 2006")

	inputDetailDuration := formatStartDate + " - " + formatEndDate

	var inputDuration string

	switch {
	case year == 1:
		inputDuration = strconv.FormatFloat(year, 'f', 0, 64) + " year"
	case year > 1:
		inputDuration = strconv.FormatFloat(year, 'f', 0, 64) + " years"
	case month == 1:
		inputDuration = strconv.FormatFloat(month, 'f', 0, 64) + " month"
	case month > 1:
		inputDuration = strconv.FormatFloat(month, 'f', 0, 64) + " months"
	case week == 1:
		inputDuration = strconv.FormatFloat(week, 'f', 0, 64) + " week"
	case week > 1:
		inputDuration = strconv.FormatFloat(week, 'f', 0, 64) + " weeks"
	case day == 1:
		inputDuration = strconv.FormatFloat(day, 'f', 0, 64) + " day"
	case day > 1:
		inputDuration = strconv.FormatFloat(day, 'f', 0, 64) + " days"
	default:
		inputDuration = "WRONG DATE!"
	}

	updateProject := Project{
		Project_name:    inputProjectName,
		Start_date:      parseStartDate,
		End_date:        parseEndDate,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
	}

	dataProject[index] = updateProject

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/home", http.StatusFound)
}
