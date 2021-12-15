package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PossibleNPC/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
	// Added in this field to the struct; deviates from code in book
	staticPath string
}

//func openDB(dsn string) (*sql.DB, error) {
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	//dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Port for Application")
	// 28NOV: Added in this functionality to serve a different directory if desired.
	serve := flag.String("serve", "./ui/static/", "Directory to serve ui files, include the trailing slash")
	// TODO: 1Dec2021 Cannot connect to the database, and I am not sure why. Internal to container works, but host
	// 			to container doesn't. Implies networking issue between host and container?
	// FIX: The solution to the above is as follows: map the 'web'@'localhost' user to 'web'@'DOCKER_IP_ADDR',
	//		then grant on the desired databases and tables, and flush privileges;
	// TODO: 12Dec2021 I dropped the database container and volumes
	// TODO: Still have to map in the user: web, to the database. There has to be a better way
	// To do this than to exec into the running container each time...
	dsn := flag.String("dsn", "web:pass@tcp(localhost)/snippetbox?parseTime=true", "MySQL data source name.")
	//dsn := flag.String("dsn", "postgres://web:example@localhost:5432/postgres", "Postgres data source name.")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		// Added in this field to the struct to make static path configurable; deviates from code in book
		staticPath: *serve,
	}

	srv := http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}