package main

import (
	"TODO/pkg/models"
	"html/template"
	"log"

	//"net"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

// Create a datastructure to hold the data
type Todo struct {
	Id   int
	Text string
}

// List of all tasks
var AllTasks []*models.ToDo

// Id for task id
var id int

// Create a handler for home
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//Send a template which will take input of creating a task and listing all the tasks
	// Parse and execute the template and pass allTasks as data
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s, err := app.todo.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	AllTasks = s
	//Render the template
	files := []string{
		"./ui/html/home.page.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}
	// app.session.Put(r, "flash", "Task Successfully created!!!!")
	err = ts.Execute(w, struct {
		Tasks []*models.ToDo
		Flash string
	}{
		Tasks: s,
		Flash: app.session.PopString(r, "flash"),
	})
	// err = ts.Execute(w, AllTasks)
	if err != nil {
		// Use the serverError() helper.
		app.serverError(w, err)
	}
}

// Create a handler for adding task
func (app *application) addTask(w http.ResponseWriter, r *http.Request) {

	id += 1 // Updating id for new task
	//Store the tasks in the structure
	task := Todo{
		Text: r.FormValue("task"),
		Id:   id,
	}
	if !app.dataValidate(r, task.Text) {
		
		//Inserting to db
		_, err := app.todo.Insert(task.Text)
		if err != nil {
			app.serverError(w, err)
			log.Println(err)
			return
		} else {
			app.session.Put(r, "flash", "Task added Successfully!!!!")

		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Create a handler for deleting task
func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	//Convert the string value and initializing it to idToDelete
	idToDelete, _ := strconv.Atoi(r.FormValue("id"))
	//log.Println(idToDelete);
	errDel := app.todo.Delete(idToDelete)
	if errDel != nil {
		log.Println(errDel)
	}
	app.session.Put(r, "flash", "Task Deleted Successfully!!!!")
	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.FormValue("id"))
	log.Println(idToUpdate)
	updateMsg := r.FormValue("update")
	app.dataValidate(r, updateMsg)
	if !app.dataValidate(r, updateMsg) {
		errUpd := app.todo.Update(idToUpdate, r.FormValue("update"))

		if errUpd != nil {
			log.Println(errUpd)
		}
		app.session.Put(r, "flash", "Task updated Successfully!!!!")
	}
	//Redirecting to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) dataValidate(r *http.Request, str string) bool {
	// Check the data validation for the users entry.
	if strings.TrimSpace(str) == "" {
		app.session.Put(r, "flash", "This field cannot be empty!!!!")
		return true
	} else if utf8.RuneCountInString(str) > 100 {
		app.session.Put(r, "flash", "This field cannot be greater than 100!!!!")
		return true
	}
	return false
}
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl", "./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}
	ts.Execute(w, app.session.PopString(r, "flash"))
	// 	app.render(w,r,"signup.page.tmpl",&templateData{
	// 		Form: forms.New(nil),
	// 	})
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	userName := r.FormValue("name")
	userEmail := r.FormValue("email")
	userPassword := r.FormValue("password")
	hashedPassword, errHashing := bcrypt.GenerateFromPassword([]byte(userPassword), 12)
	if errHashing != nil {
		log.Println(errHashing)
		return
	}

	err := app.users.Insert(userName, userEmail, hashedPassword)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.session.Put(r, "flash", "User already exist")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
	} else {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "flash", "Sign up successfull!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl", "./ui/html/base.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}
	ts.Execute(w, app.session.PopString(r, "flash"))

}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("email")
	userPassword := r.FormValue("password")
	isUser, err := app.users.Authenticate(userName, userPassword)
	log.Println(isUser)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	if isUser != -1 {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "flash", " Login Successfull!")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		app.session.Put(r, "flash", " Login Failed!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.session.Put(r, "Authenticated", false)
	}
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Put(r, "Authenticated", false)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	app.session.Put(r, "flash", " Successfully Logged Out")
}

func (app *application) spcl_taskAdd(w http.ResponseWriter, r *http.Request) {
	id += 1
	task := Todo{
		Text: r.FormValue("task"),
		Id:   id,
	}
	if !app.dataValidate(r, task.Text) {
		_, err := app.spcl_todo.SpclTaskAdd(task.Text)
		if err != nil {
			app.serverError(w, err)
			log.Println(err)
			return
		} else {
			app.session.Put(r, "flash", "Task added Successfully!!!!")
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) spcl_taskDelete(w http.ResponseWriter, r *http.Request) {
	//Convert the string value and initializing it to idToDelete
	idToDelete, _ := strconv.Atoi(r.FormValue("id"))
	//log.Println(idToDelete);
	errDel := app.spcl_todo.SpclTaskDelete(idToDelete)
	if errDel != nil {
		log.Println(errDel)
	}
	app.session.Put(r, "flash", "Task Deleted Successfully!!!!")
}
