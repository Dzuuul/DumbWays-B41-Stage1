package main

import (
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
	Project_name    string
	Start_date      string
	End_date        string
	Detail_duration string
	Duration        string
	Description     string
	React_js        string
	Vue_js          string
	Angular         string
	Laravel         string
	Icon_react_js   string
	Icon_vue_js     string
	Icon_angular    string
	Icon_laravel    string
}

var dataProject = []Project{}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/my-project", project).Methods("GET")
	r.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	r.HandleFunc("/form-edit-project/{index}", formEditProject).Methods("GET")
	r.HandleFunc("/detail-project/{index}", detailProject).Methods("GET")
	r.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	r.HandleFunc("/add-my-project", addProject).Methods("POST")
	r.HandleFunc("/edit-my-project/{index}", editProject).Methods("POST")

	fmt.Println("Server is running on port 5656...\t(press \"ctrl + c\" to abort)")
	http.ListenAndServe("localhost:5656", r)
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

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/my-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Projects": dataProject,
	}

	tmpl.Execute(w, resp)
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
				React_js:        data.React_js,
				Vue_js:          data.Vue_js,
				Angular:         data.Angular,
				Laravel:         data.Laravel,
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

	rjs := r.PostForm.Get("check-reactjs")
	vjs := r.PostForm.Get("check-vuejs")
	ang := r.PostForm.Get("check-angular")
	lar := r.PostForm.Get("check-laravel")

	detailReactJs, cardReactJs := rjs, rjs
	detailVueJs, cardVueJs := vjs, vjs
	detailAngular, cardAngular := ang, ang
	detailLaravel, cardLaravel := lar, lar

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

	switch detailReactJs {
	case "on":
		detailReactJs = `<p><i class="fa-brands fa-react fa-xl me-2"></i>React Js</p>`
	default:
		detailReactJs = ""
	}
	switch detailVueJs {
	case "on":
		detailVueJs = `<p><i class="fa-brands fa-vuejs fa-xl me-2"></i>Vue Js</p>`
	default:
		detailVueJs = ""
	}
	switch detailAngular {
	case "on":
		detailAngular = `<p><i class="fa-brands fa-angular fa-xl me-2"></i>Angular</p>`
	default:
		detailAngular = ""
	}
	switch detailLaravel {
	case "on":
		detailLaravel = `<p><i class="fa-brands fa-laravel fa-xl me-2"></i>Laravel</p>`
	default:
		detailLaravel = ""
	}

	switch cardReactJs {
	case "on":
		cardReactJs = `<i class="fa-brands fa-react fa-xl me-2"></i>`
	default:
		cardReactJs = ""
	}
	switch cardVueJs {
	case "on":
		cardVueJs = `<i class="fa-brands fa-vuejs fa-xl me-2"></i>`
	default:
		cardVueJs = ""
	}
	switch cardAngular {
	case "on":
		cardAngular = `<i class="fa-brands fa-angular fa-xl me-2"></i>`
	default:
		cardAngular = ""
	}
	switch cardLaravel {
	case "on":
		cardLaravel = `<i class="fa-brands fa-laravel fa-xl me-2"></i>`
	default:
		cardLaravel = ""
	}

	newProject := Project{
		Project_name:    inputProjectName,
		Start_date:      inputStartDate,
		End_date:        inputEndDate,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
		React_js:        detailReactJs,
		Vue_js:          detailVueJs,
		Angular:         detailAngular,
		Laravel:         detailLaravel,
		Icon_react_js:   cardReactJs,
		Icon_vue_js:     cardVueJs,
		Icon_angular:    cardAngular,
		Icon_laravel:    cardLaravel,
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/my-project", http.StatusMovedPermanently)
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
				React_js:        data.React_js,
				Vue_js:          data.Vue_js,
				Angular:         data.Angular,
				Laravel:         data.Laravel,
			}
		}
	}

	data := map[string]interface{}{
		"Index": index,
		"Edit":  edit,
	}

	tmpl.Execute(w, data)
	http.Redirect(w, r, "/my-project", http.StatusMovedPermanently)
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

	rjs := r.PostForm.Get("check-reactjs")
	vjs := r.PostForm.Get("check-vuejs")
	ang := r.PostForm.Get("check-angular")
	lar := r.PostForm.Get("check-laravel")

	detailReactJs, cardReactJs := rjs, rjs
	detailVueJs, cardVueJs := vjs, vjs
	detailAngular, cardAngular := ang, ang
	detailLaravel, cardLaravel := lar, lar

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

	switch detailReactJs {
	case "on":
		detailReactJs = `<p><i class="fa-brands fa-react fa-xl me-2"></i>React Js</p>`
	default:
		detailReactJs = ""
	}
	switch detailVueJs {
	case "on":
		detailVueJs = `<p><i class="fa-brands fa-vuejs fa-xl me-2"></i>Vue Js</p>`
	default:
		detailVueJs = ""
	}
	switch detailAngular {
	case "on":
		detailAngular = `<p><i class="fa-brands fa-angular fa-xl me-2"></i>Angular</p>`
	default:
		detailAngular = ""
	}
	switch detailLaravel {
	case "on":
		detailLaravel = `<p><i class="fa-brands fa-laravel fa-xl me-2"></i>Laravel</p>`
	default:
		detailLaravel = ""
	}

	switch cardReactJs {
	case "on":
		cardReactJs = `<i class="fa-brands fa-react fa-xl me-2"></i>`
	default:
		cardReactJs = ""
	}
	switch cardVueJs {
	case "on":
		cardVueJs = `<i class="fa-brands fa-vuejs fa-xl me-2"></i>`
	default:
		cardVueJs = ""
	}
	switch cardAngular {
	case "on":
		cardAngular = `<i class="fa-brands fa-angular fa-xl me-2"></i>`
	default:
		cardAngular = ""
	}
	switch cardLaravel {
	case "on":
		cardLaravel = `<i class="fa-brands fa-laravel fa-xl me-2"></i>`
	default:
		cardLaravel = ""
	}

	updateProject := Project{
		Project_name:    inputProjectName,
		Start_date:      inputStartDate,
		End_date:        inputEndDate,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
		React_js:        detailReactJs,
		Vue_js:          detailVueJs,
		Angular:         detailAngular,
		Laravel:         detailLaravel,
		Icon_react_js:   cardReactJs,
		Icon_vue_js:     cardVueJs,
		Icon_angular:    cardAngular,
		Icon_laravel:    cardLaravel,
	}

	dataProject[index] = updateProject

	http.Redirect(w, r, "/my-project", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/my-project", http.StatusFound)
}
