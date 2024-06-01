package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"snipit.bikraj.net/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Http Network Address")

	dsn := flag.String(
		"dsn",
		" ",
		"MySQL Data Source name",
	)
	// Custom Loggers

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	flag.Parse()
	fmt.Println(*dsn)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache, err := newTemplateCache()
	defer db.Close()
	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	mux := app.routes()
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Println("Starting Server on", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

//	func neuter(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if strings.HasSuffix(r.URL.Path, "/") {
//				http.NotFound(w, r)
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

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
