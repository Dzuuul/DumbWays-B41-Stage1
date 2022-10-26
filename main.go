package main

import (
	"context"
	"day-10/connection"
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

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/home", project).Methods("GET")
	r.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	r.HandleFunc("/form-edit-project/{id}", formEditProject).Methods("GET")
	r.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")
	r.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	r.HandleFunc("/add-my-project", addProject).Methods("POST")
	r.HandleFunc("/edit-my-project/{id}", editProject).Methods("POST")

	fmt.Println("Server is running on port 5656...\t(press \"ctrl + c\" to abort)")
	http.ListenAndServe("localhost:5656", r)
}

func home(w http.ResponseWriter, r *http.Request) {
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

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects ORDER BY id DESC")

	var result []Project
	for data.Next() {
		var each = Project{}

		err = data.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		hour := each.End_date.Sub(each.Start_date).Hours()
		day := hour / 24
		week := day / 7
		month := day / 30
		year := day / 365

		var inputDuration string

		switch {
		case day == 1:
			inputDuration = strconv.Itoa(int(day)) + " day"
		case day > 1 && day <= 6:
			inputDuration = strconv.Itoa(int(day)) + " days"
		case day == 7:
			inputDuration = strconv.Itoa(int(week)) + " week"
		case day > 7 && day <= 29:
			inputDuration = strconv.Itoa(int(week)) + " weeks"
		case day == 30:
			inputDuration = strconv.Itoa(int(month)) + " month"
		case day > 30 && day <= 364:
			inputDuration = strconv.Itoa(int(month)) + " months"
		case day == 365:
			inputDuration = strconv.Itoa(int(year)) + " year"
		case day > 365:
			inputDuration = strconv.Itoa(int(year)) + " years"
		default:
			inputDuration = "WRONG DATE!"
		}

		each.Duration = inputDuration

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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/detail-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	ProjectDetail := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Project_name, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	hour := ProjectDetail.End_date.Sub(ProjectDetail.Start_date).Hours()
	day := hour / 24
	week := day / 7
	month := day / 30
	year := day / 365

	formatStartDate := ProjectDetail.Start_date.Format("2 Jan 2006")
	formatEndDate := ProjectDetail.End_date.Format("2 Jan 2006")

	inputDetailDuration := formatStartDate + " - " + formatEndDate

	var inputDuration string

	switch {
	case day == 1:
		inputDuration = strconv.Itoa(int(day)) + " day"
	case day > 1 && day <= 6:
		inputDuration = strconv.Itoa(int(day)) + " days"
	case day == 7:
		inputDuration = strconv.Itoa(int(week)) + " week"
	case day > 7 && day <= 29:
		inputDuration = strconv.Itoa(int(week)) + " weeks"
	case day == 30:
		inputDuration = strconv.Itoa(int(month)) + " month"
	case day > 30 && day <= 364:
		inputDuration = strconv.Itoa(int(month)) + " months"
	case day == 365:
		inputDuration = strconv.Itoa(int(year)) + " year"
	case day > 365:
		inputDuration = strconv.Itoa(int(year)) + " years"
	default:
		inputDuration = "WRONG DATE!"
	}

	ProjectDetail.Detail_duration = inputDetailDuration
	ProjectDetail.Duration = inputDuration

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

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(name, start_date, end_date, description, technologies) VALUES ($1, $2, $3, $4, $5)", inputProjectName, inputStartDate, inputEndDate, inputDescription, inputTechnologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/edit-my-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	UpdateProject := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies FROM tb_projects WHERE id=$1", id).Scan(&UpdateProject.Id, &UpdateProject.Project_name, &UpdateProject.Start_date, &UpdateProject.End_date, &UpdateProject.Description, &UpdateProject.Technologies)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	hour := UpdateProject.End_date.Sub(UpdateProject.Start_date).Hours()
	day := hour / 24
	week := day / 7
	month := day / 30
	year := day / 365

	formatStartDate := UpdateProject.Start_date.Format("2 Jan 2006")
	formatEndDate := UpdateProject.End_date.Format("2 Jan 2006")

	inputDetailDuration := formatStartDate + " - " + formatEndDate

	var inputDuration string

	switch {
	case day == 1:
		inputDuration = strconv.Itoa(int(day)) + " day"
	case day > 1 && day <= 6:
		inputDuration = strconv.Itoa(int(day)) + " days"
	case day == 7:
		inputDuration = strconv.Itoa(int(week)) + " week"
	case day > 7 && day <= 29:
		inputDuration = strconv.Itoa(int(week)) + " weeks"
	case day == 30:
		inputDuration = strconv.Itoa(int(month)) + " month"
	case day > 30 && day <= 364:
		inputDuration = strconv.Itoa(int(month)) + " months"
	case day == 365:
		inputDuration = strconv.Itoa(int(year)) + " year"
	case day > 365:
		inputDuration = strconv.Itoa(int(year)) + " years"
	default:
		inputDuration = "WRONG DATE!"
	}

	UpdateProject.Detail_duration = inputDetailDuration
	UpdateProject.Duration = inputDuration

	// for i, data := range dataProject {
	// 	if id == i {
	// 		UpdateProject = Project{
	// 			Id:              data.Id,
	// 			Project_name:    data.Project_name,
	// 			Start_date:      data.Start_date,
	// 			End_date:        data.End_date,
	// 			Detail_duration: data.Detail_duration,
	// 			Duration:        data.Duration,
	// 			Description:     data.Description,
	// 		}
	// 	}
	// }

	data := map[string]interface{}{
		"Update": UpdateProject,
	}

	// edit := Project{}

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// data := map[string]interface{}{
	// 	"Edit": edit,
	// }

	tmpl.Execute(w, data)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
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
	case day == 1:
		inputDuration = strconv.Itoa(int(day)) + " day"
	case day > 1 && day <= 6:
		inputDuration = strconv.Itoa(int(day)) + " days"
	case day == 7:
		inputDuration = strconv.Itoa(int(week)) + " week"
	case day > 7 && day <= 29:
		inputDuration = strconv.Itoa(int(week)) + " weeks"
	case day == 30:
		inputDuration = strconv.Itoa(int(month)) + " month"
	case day > 30 && day <= 364:
		inputDuration = strconv.Itoa(int(month)) + " months"
	case day == 365:
		inputDuration = strconv.Itoa(int(year)) + " year"
	case day > 365:
		inputDuration = strconv.Itoa(int(year)) + " years"
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

	dataProject[id] = updateProject

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/home", http.StatusFound)
}
