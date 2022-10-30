package main

import (
	"context"
	"day-12/connection"
	"day-12/middleware"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title, UserName, FlashData string
	IsLogin                    bool
}

var Data = MetaData{}

type Project struct {
	Id           int
	User_id      int
	IsLogin      bool
	Technologies []string
	Start_date   time.Time
	End_date     time.Time
	Project_name, Start_date_string,
	End_date_string, Detail_duration,
	Duration, Description, Image string
}

type User struct {
	Name, Email, Password string
	Id                    int
}

func main() {
	r := mux.NewRouter()

	connection.DatabaseConnect()

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/home", project).Methods("GET")
	r.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	r.HandleFunc("/form-edit-project/{id}", formEditProject).Methods("GET")
	r.HandleFunc("/detail-project/{id}", detailProject).Methods("GET")
	r.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	r.HandleFunc("/add-my-project", middleware.UploadFile(addProject)).Methods("POST")
	r.HandleFunc("/edit-my-project/{id}", middleware.UploadFile(editProject)).Methods("POST")

	r.HandleFunc("/form-register", formRegister).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")

	r.HandleFunc("/form-login", formLogin).Methods("GET")
	r.HandleFunc("/login", login).Methods("POST")

	r.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server is running on port 5656...\t(press \"ctrl + c\" to exit)")
	http.ListenAndServe("localhost:5656", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/home.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")

	if session.Values["IsLogin"] != true {

		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, tb_projects.name, start_date, end_date, description, technologies, image FROM tb_projects ORDER BY id DESC")

		var result []Project
		for data.Next() {
			var each = Project{}

			err = data.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)

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
			"Data":     Data,
			"Projects": result,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)

	} else {

		userID := session.Values["Id"].(int)

		data, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, tb_projects.name, start_date, end_date, description, technologies, image FROM tb_projects WHERE tb_projects.user_id = $1 ORDER BY id DESC", userID)

		var result []Project
		for data.Next() {
			var each = Project{}

			err = data.Scan(&each.Id, &each.Project_name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)

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
			"Data":     Data,
			"Projects": result,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)

	}

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

	var Data = MetaData{}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	ProjectDetail := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Project_name, &ProjectDetail.Start_date, &ProjectDetail.End_date, &ProjectDetail.Description, &ProjectDetail.Technologies, &ProjectDetail.Image)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
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
		"Data":    Data,
		"Project": ProjectDetail,
	}

	w.WriteHeader(http.StatusOK)
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

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	tmpl.Execute(w, Data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	inputUserId := session.Values["Id"].(int)

	inputProjectName := r.PostForm.Get("input-project-name")
	inputStartDate := r.PostForm.Get("input-start-date")
	inputEndDate := r.PostForm.Get("input-end-date")
	inputDescription := r.PostForm.Get("input-description")

	var inputTechnologies []string
	inputTechnologies = r.Form["technologies"]

	dataContex := r.Context().Value("dataFile")
	inputImage := dataContex.(string)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(user_id, name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6, $7)", inputUserId, inputProjectName, inputStartDate, inputEndDate, inputDescription, inputTechnologies, inputImage)

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

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}
	tmpl.Execute(w, Data)
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

	var Data = MetaData{}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["Name"].(string)
	}

	UpdateProject := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&UpdateProject.Id, &UpdateProject.Project_name, &UpdateProject.Start_date, &UpdateProject.End_date, &UpdateProject.Description, &UpdateProject.Technologies, &UpdateProject.Image)

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

	UpdateProject.Start_date_string = UpdateProject.Start_date.Format("2006-01-02")
	UpdateProject.End_date_string = UpdateProject.End_date.Format("2006-01-02")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
	}

	data := map[string]interface{}{
		"Data":   Data,
		"Update": UpdateProject,
	}

	tmpl.Execute(w, data)
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

	dataContex := r.Context().Value("dataFile")
	inputImage := dataContex.(string)

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET name = $1, start_date = $2, end_date = $3, description = $4, technologies = $5, image = $6 WHERE id = $7", inputProjectName, inputStartDate, inputEndDate, inputDescription, inputTechnologies, inputImage, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error message : " + err.Error()))
		return
	}

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

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var email = r.PostForm.Get("inputEmail")
	var password = r.PostForm.Get("inputPassword")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_users (name, email, password) VALUES($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.AddFlash("Registration Successful!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	fm := session.Flashes("message")

	var flashes []string
	if len(fm) > 0 {
		session.Save(r, w)
		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	Data.FlashData = strings.Join(flashes, "")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func login(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("inputEmail")
	password := r.PostForm.Get("inputPassword")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_users WHERE email=$1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Email is not registered!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		session.AddFlash("Wrong password!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/form-login", http.StatusMovedPermanently)
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Email"] = user.Email
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800

	session.AddFlash("Successfully Login!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout!")
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
