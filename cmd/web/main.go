package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Prateek2593/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// define an application struct to hold the application wide dependencies for the web application. for now we'll only include fields for the two custom loggers, but we'll add more to it as the build progress
// add a templateCache field to application struct
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// add a snippets field to the application struct. this will allow us to make the SnippetModel object available to our handlers
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder // add a formDecoder field to hold a pointer to a form.Decoder instance
}

func main() {

	// define a new command line flag with name addr, a default value of ":4000" and some short help text explaining what the flag controls. the value of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// define a new command line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// importantly, we use the flag.Parse() function to parse the command line flag. this reads in command line flag value and assigns it to the addr variable. you need to call this *before* you use the addr variable otherwise it will always contain the default value of ":4000". if any errors are encountered during parsing the application will be terminated
	flag.Parse()

	// use log.New() to create a logger for writing information messages. this takes three parameters: the destination to write the logs to(os.Stdout), a string prefix for message (INFO followed by tab), and flags to indicate what additional information to include(local date and time). note that the flags are joined using the bitwise OR operator
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// create a logger for writing error messages in the same way, but use stderr as the destination and use the log.Lshortfile flag to include the relevant file name and line number
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// to keep the main() function tidy we have put the code for creating a connection pool into separate openDB() function below. we pass openDB() tht dsn from command line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// we also defer a call to db.Close(), so that the connection pool is closed before main function exits
	defer db.Close()

	// initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialize a decoder instance
	formDecoder := form.NewDecoder()

	// create a new instance of our application struct with the custom loggers
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// initialize a models.SnippetModel instance and add it to the application dependencies
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache, // add it to application dependencies
		formDecoder:   formDecoder,   // add it to application dependencies,
	}

	// initialize a new http.Server struct. we set the addr and handler fields so that the server uses the same network address and routes as before
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,     // set the error log to the custom error logger we created earlier
		Handler:  app.routes(), // call the new app.routes() method to get the servemux containing our routes
	}

	//the value returned from the flag.String() function is a pointer to the flag value, not the value itself. so we need to dereference the pointer(i.e. prefix it with the *symbol) before using it. note that we're using the log.Printf() function to interpolate the address with log message.
	infoLog.Printf("Starting server on %s", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in two parameters: TCP network address to listen on(in this case":4000) and the servemux we just created. If http.ListenAndServe() return an error, we use the log.Fatal() function to log the error message and exit. note that any error returned by http.ListenAndServe() is always non-nil
	// err := http.ListenAndServe(*addr, mux)

	// call the ListenAndServe() method on our new http.Server struct
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// the openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN
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
