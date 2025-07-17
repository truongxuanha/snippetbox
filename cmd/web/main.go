package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.xtruong27.net/interval/models"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateData map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dns", "web:pass@/snippetbox?parseTime=true", "MySql data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db,err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
		
	}
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateData: templateCache,
	}
	
	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	
	err = srv.ListenAndServe()

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping();err != nil {
		return nil, err
	}
	return db, nil
}