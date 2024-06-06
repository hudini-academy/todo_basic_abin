package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// It sends a 404 Not Found response the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
func initLoggers() (*log.Logger, *log.Logger) {
	// Create a file destination predefined in the main that holds the info
	f, error := os.OpenFile("./tmp/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if error != nil {
		log.Fatal(error)
	}


	// We create a new custom
	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	// Create a file destination predefined in the main that holds the error logs
	Er, error := os.OpenFile("./tmp/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if error != nil {
		log.Fatal(error)
	}
	// Create a logger for writing error messages in the same way,
	errorLog := log.New(Er, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}
