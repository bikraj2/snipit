package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
  _ "github.com/go-sql-driver/mysql"
	"snipit.bikraj.net/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets       models.SnippetModelInterface 
  users          models.UserModelInterface
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
  sessionManager *scs.SessionManager
  debug  bool
}

func main() {
	addr := flag.String("addr", ":4000", "Http Network Address")
  debug := flag.Bool("debug", false,  "Run the applicaiton in Debug Mode")

	dsn := flag.String(
		"dsn",
	"",	
		"MySQL Data Source name",
	)
	// Custom Loggers

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	flag.Parse()
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	templateCache, err := newTemplateCache()
  if err!=nil {
  panic(err)
  }

	formDecoder := form.NewDecoder()
  sessionManager:=scs.New()
  sessionManager.Store = mysqlstore.New(db)
  sessionManager.Lifetime = 12 * time.Hour
	defer db.Close()
	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
    users:       &models.UserModel{Db: db}, 
		templateCache: templateCache,
		formDecoder:   formDecoder,
    sessionManager: sessionManager,
    debug : *debug,
	}

  tlsConfig := &tls.Config {
    CurvePreferences: []tls.CurveID{tls.X25519,tls.CurveP256},
  CipherSuites: []uint16{
      tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
      tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, 
      tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, 
      tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, 
      tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, 
      tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
},
  } 
	// mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	mux := app.routes()
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
    TLSConfig: tlsConfig,
    IdleTimeout: time.Minute,
    ReadTimeout: 5* time.Second,
    WriteTimeout: 10* time.Second,
	}

	infoLog.Println("Starting Server on", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem","./tls/key.pem")
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
