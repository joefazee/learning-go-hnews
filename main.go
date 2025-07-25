package main

import (
	"database/sql"
	"github.com/golangcollege/sessions"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	userRepo    UserRepository
	postRepo    PostRepository
	templateDir string
	publicPath  string
	tp          *TemplateRenderer
	session     *sessions.Session
}

func main() {

	db, err := connectToDatabase("users_database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte("u46IpCV9y5VlurXXXODJEhgOY8m9JVE4"))
	session.Lifetime = 24 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteLaxMode

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ltime|log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ltime|log.LstdFlags),
		userRepo:    NewSQLUserRepository(db),
		postRepo:    NewSQLPostRepository(db),
		templateDir: "./templates",
		publicPath:  "./public",
		session:     session,
	}
	app.tp = NewTemplateRenderer(app.templateDir, false)

	log.Println("Listening on :8080")
	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

}

func connectToDatabase(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
