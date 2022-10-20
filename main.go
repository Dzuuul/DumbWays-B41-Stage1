package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

type Project struct {
	Project_name    string
	Detail_duration string
	Duration        string
	Description     string
	React_js        string
	Vue_js          string
	Angular         string
	Laravel         string
}

var dataProject = []Project{}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/detail-project/{i}", detailProject).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/add-project", projectPage).Methods("GET")
	r.HandleFunc("/update-project", updateProjectPage).Methods("GET")
	r.HandleFunc("/my-project", submitProject).Methods("POST")
	r.HandleFunc("/edit-project", editProject).Methods("POST")
	r.HandleFunc("/delete-project", deleteProject).Methods("GET")

	fmt.Println("Server is running on port 5678...\t(press \"ctrl + c\" to cancel)")
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

func projectPage(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	ProjectDetail := Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["i"])

	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				Project_name: data.Project_name,
				Duration:     data.Duration,
				Description:  data.Description,
				React_js:     data.React_js,
				Vue_js:       data.Vue_js,
				Angular:      data.Angular,
				Laravel:      data.Laravel,
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

	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	tmpl.Execute(w, nil)
}

func submitProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	inputProjectName := r.PostForm.Get("input-project-name")
	inputStartDate := r.PostForm.Get("input-start-date")
	inputEndDate := r.PostForm.Get("input-end-date")
	inputDescription := r.PostForm.Get("input-description")

	checkReactJs := r.PostForm.Get("check-reactjs")
	checkVueJs := r.PostForm.Get("check-vuejs")
	checkAngular := r.PostForm.Get("check-angular")
	checkLaravel := r.PostForm.Get("check-laravel")

	timeStartDate, _ := time.Parse("2006-01-02", inputStartDate)
	timeEndDate, _ := time.Parse("2006-01-02", inputEndDate)

	formatStartDate := timeStartDate.Format("2 Jan 2006")
	formatEndDate := timeEndDate.Format("2 Jan 2006")

	dateDifference := timeEndDate.Sub(timeStartDate)
	dayDuration := int64(dateDifference.Hours() / 24)

	inputDetailDuration := formatStartDate + " - " + formatEndDate
	inputDuration := duration(int(dayDuration))

	fmt.Printf("\nCalculate Duration\t:= %v\n", inputDetailDuration)
	fmt.Printf("\nDate Difference\t\t:= %v\n", dateDifference)

	if dayDuration == 1 {
		fmt.Printf("\nDay Duration\t\t= %v day\n", dayDuration)
	} else {
		fmt.Printf("\nDay Duration\t\t= %v days\n", dayDuration)
	}

	fmt.Printf("\nDuration\t\t= %v\n", inputDuration)

	switch checkReactJs {
	case "on":
		checkReactJs = `<p><i class="fa-brands fa-react fa-xl me-2"></i>React Js</p>`
	default:
		checkReactJs = ""
	}
	switch checkVueJs {
	case "on":
		checkVueJs = `<p><i class="fa-brands fa-vuejs fa-xl me-2"></i>Vue Js</p>`
	default:
		checkVueJs = ""
	}
	switch checkAngular {
	case "on":
		checkAngular = `<p><i class="fa-brands fa-angular fa-xl me-2"></i>Angular</p>`
	default:
		checkAngular = ""
	}
	switch checkLaravel {
	case "on":
		checkLaravel = `<p><i class="fa-brands fa-laravel fa-xl me-2"></i>Laravel</p>`
	default:
		checkLaravel = ""
	}

	fmt.Println("______")
	fmt.Println("Form Result :")
	fmt.Printf("\nProject Name\t: %v\n\nDuration\t: %v\n\nDescription\t:\n%v\n\n", inputProjectName, inputDuration, inputDescription)
	fmt.Println("Technologies\t:")
	if checkReactJs != "" {
		fmt.Printf("  ✔ React Js ✔ ")
	}
	if checkVueJs != "" {
		fmt.Printf("  ✔ Vue Js ✔ ")
	}
	if checkAngular != "" {
		fmt.Printf("  ✔ Angular ✔ ")
	}
	if checkLaravel != "" {
		fmt.Printf("  ✔ Laravel  ✔ ")
	}
	fmt.Println("\n\n______")

	addProject := Project{
		Project_name:    inputProjectName,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
		React_js:        checkReactJs,
		Vue_js:          checkVueJs,
		Angular:         checkAngular,
		Laravel:         checkLaravel,
	}

	dataProject = append(dataProject, addProject)

	http.Redirect(w, r, "/my-project", http.StatusMovedPermanently)
}

func updateProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
	}

	i, _ := strconv.Atoi(mux.Vars(r)["i"])

	data := map[string]interface{}{
		"id": i,
	}

	tmpl.Execute(w, data)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	i, _ := strconv.Atoi(mux.Vars(r)["i"])
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	inputProjectName := r.PostForm.Get("input-project-name")
	inputStartDate := r.PostForm.Get("input-start-date")
	inputEndDate := r.PostForm.Get("input-end-date")
	inputDescription := r.PostForm.Get("input-description")

	checkReactJs := r.PostForm.Get("check-reactjs")
	checkVueJs := r.PostForm.Get("check-vuejs")
	checkAngular := r.PostForm.Get("check-angular")
	checkLaravel := r.PostForm.Get("check-laravel")

	timeStartDate, _ := time.Parse("2006-01-02", inputStartDate)
	timeEndDate, _ := time.Parse("2006-01-02", inputEndDate)

	formatStartDate := timeStartDate.Format("2 Jan 2006")
	formatEndDate := timeEndDate.Format("2 Jan 2006")

	dateDifference := timeEndDate.Sub(timeStartDate)
	dayDuration := int64(dateDifference.Hours() / 24)

	inputDetailDuration := formatStartDate + " - " + formatEndDate
	inputDuration := duration(int(dayDuration))

	fmt.Printf("\nDetail Duration\t:= %v\n", inputDetailDuration)
	fmt.Printf("\nDate Difference\t\t:= %v\n", dateDifference)

	if dayDuration == 1 {
		fmt.Printf("\nDay Duration\t\t= %v day\n", dayDuration)
	} else {
		fmt.Printf("\nDay Duration\t\t= %v days\n", dayDuration)
	}

	fmt.Printf("\nDuration\t\t= %v\n", inputDuration)

	switch checkReactJs {
	case "on":
		checkReactJs = `<p><i class="fa-brands fa-react fa-xl me-2"></i>React Js</p>`
	default:
		checkReactJs = ""
	}
	switch checkVueJs {
	case "on":
		checkVueJs = `<p><i class="fa-brands fa-vuejs fa-xl me-2"></i>Vue Js</p>`
	default:
		checkVueJs = ""
	}
	switch checkAngular {
	case "on":
		checkAngular = `<p><i class="fa-brands fa-angular fa-xl me-2"></i>Angular</p>`
	default:
		checkAngular = ""
	}
	switch checkLaravel {
	case "on":
		checkLaravel = `<p><i class="fa-brands fa-laravel fa-xl me-2"></i>Laravel</p>`
	default:
		checkLaravel = ""
	}

	fmt.Println("______")
	fmt.Println("Form Result :")
	fmt.Printf("\nProject Name\t: %v\n\nDuration\t: %v\n\nDescription\t:\n%v\n\n", inputProjectName, inputDuration, inputDescription)
	fmt.Println("Technologies\t:")
	if checkReactJs != "" {
		fmt.Printf("  ✔ React Js ✔ ")
	}
	if checkVueJs != "" {
		fmt.Printf("  ✔ Vue Js ✔ ")
	}
	if checkAngular != "" {
		fmt.Printf("  ✔ Angular ✔ ")
	}
	if checkLaravel != "" {
		fmt.Printf("  ✔ Laravel  ✔ ")
	}
	fmt.Println("\n\n______")

	updateProject := Project{
		Project_name:    inputProjectName,
		Detail_duration: inputDetailDuration,
		Duration:        inputDuration,
		Description:     inputDescription,
		React_js:        checkReactJs,
		Vue_js:          checkVueJs,
		Angular:         checkAngular,
		Laravel:         checkLaravel,
	}

	dataProject[i] = updateProject

	http.Redirect(w, r, "/my-project", http.StatusMovedPermanently)

}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	i, _ := strconv.Atoi(mux.Vars(r)["i"])

	dataProject = append(dataProject[:i], dataProject[i+1:]...)

	http.Redirect(w, r, "/my-project", http.StatusFound)
}

func duration(d int) string {
	dItoa := strconv.Itoa(int(d))

	if d == 1 || d == 0 {
		return dItoa + " day"
	} else if d < 0 {
		return "-"
	} else if d > 1 && d < 7 {
		return dItoa + " days"
	} else if d == 7 {
		weekD := d / 7
		weekDItoa := strconv.Itoa(int(weekD))
		return weekDItoa + " week"
	} else if d > 7 && d < 30 {
		weekD := d / 7
		weekDItoa := strconv.Itoa(int(weekD))
		return weekDItoa + " weeks"
	} else if d == 30 {
		monthD := d / 30
		monthDItoa := strconv.Itoa(int(monthD))
		return monthDItoa + " month"
	} else if d > 30 && d < 365 {
		monthD := d / 30
		monthDItoa := strconv.Itoa(int(monthD))
		return monthDItoa + " months"
	} else if d == 365 {
		yearD := d / 365
		yearDItoa := strconv.Itoa(int(yearD))
		return yearDItoa + " year"
	} else if d > 365 {
		yearD := d / 365
		yearDItoa := strconv.Itoa(int(yearD))
		return yearDItoa + " years"
	}
	return dItoa
}
