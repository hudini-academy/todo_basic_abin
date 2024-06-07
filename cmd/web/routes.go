package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	// Checking the authentication middleware
	dynamicMiddleware1 := alice.New(app.session.Enable, app.middlewareAuthenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware1.ThenFunc(app.home))
	mux.Post("/addTaskPage", dynamicMiddleware1.ThenFunc(app.addTask))
	mux.Post("/tasksdelete", dynamicMiddleware1.ThenFunc(app.deleteTask))
	mux.Post("/updatetask", dynamicMiddleware1.ThenFunc(app.updateTask))
	mux.Post("/spclAddtask",dynamicMiddleware.ThenFunc(app.spcl_taskAdd))
	mux.Post("/spclDeletask",dynamicMiddleware.ThenFunc(app.spcl_taskDelete))

	//User authentication
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware1.ThenFunc(app.logoutUser))
	// Have  a fileserver from which CSS file will be served
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	return standardMiddleware.Then(mux)
}
