package main

import (
	"TODO/pkg/models/mysql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// Creating a type config to declare custom typed server address and other variables in the command line
type Config struct {
	Addr      string
	StaticDir string
}

// Define an application struct to hold the application-wide dependencies for the web application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   *Config
	todo     *mysql.ToDoModel
	session  *sessions.Session
	users    *mysql.UserModel
	spcl_todo *mysql.SpclToDoModel
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static directory")
	flag.Parse()

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "root:root@/todo?parseTime=true", "MySQL database")

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")
	flag.Parse()

	// Initializing infoLog and errorLog from the helpers.go
	infoLog, errorLog := initLoggers()
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// The code for creating a connection pool into the separate openDB() function below.
	// We pass openDB() the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed before the main() function exits.
	defer db.Close()

	// Creating an instance named 'app' to call the struct
	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		config:    cfg,
		session:   session,
		todo:      &mysql.ToDoModel{DB: db},
		users:     &mysql.UserModel{DB: db},
		spcl_todo: &mysql.SpclToDoModel{DB: db},
	}

	// Rewrite the source code with our own server address, handler and error logs.
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	// Call the ListenAndServe() method on our new http.Server struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
