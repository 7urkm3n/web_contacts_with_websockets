package main

import (
	"backend/internal/models"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type config struct {
	port int
	env  string
	db   struct {
		reset bool
	}
}

type application struct {
	config config
	logger *slog.Logger
	models models.Models
	ws     *websocket.Conn
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.BoolVar(&cfg.db.reset, "resetDB", false, "Sqlite3 database reset default false")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openSqlite(&cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())

	go app.pushContact()

	os.Exit(1)
}

func openSqlite(cfg *config) (*sql.DB, error) {
	if cfg.db.reset {
		log.Println("Resetting database...")
		os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	}

	var err error
	db, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err)
	}

	// SQL statement to create the contacts table if it doesn't exist
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS contacts (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL UNIQUE,
		phone_number TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS contact_edit_history (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		contact_id INTEGER,
		changes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (contact_id)
			REFERENCES contacts (id) 
					ON DELETE CASCADE 
					ON UPDATE NO ACTION
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, sqlStmt)
	}

	return db, nil
}
