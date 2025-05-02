package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Default schema to use if no schema.sql file is found
const defaultSchema = `
-- Routes table to define HTTP endpoints
CREATE TABLE IF NOT EXISTS routes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    handler_function TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(path, method)
);

-- Templates table for storing HTML templates
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Static files table for storing static assets
CREATE TABLE IF NOT EXISTS static_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL UNIQUE,
    content BLOB NOT NULL,
    mime_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Functions table for storing SQL functions
CREATE TABLE IF NOT EXISTS functions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    sql_code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func initDB(dbPath string, reset bool) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if reset {
		// Drop all tables if they exist
		tables := []string{"routes", "templates", "static_files", "functions"}
		for _, table := range tables {
			_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
			if err != nil {
				return fmt.Errorf("failed to drop table %s: %v", table, err)
			}
		}
		log.Println("Database reset complete")
	}

	// Try to read schema from file first
	schema, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		// If schema.sql doesn't exist, use default schema
		schema = []byte(defaultSchema)
	}

	// Execute schema
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Find matching route
	var handlerFunc string
	err := db.QueryRow("SELECT handler_function FROM routes WHERE path = ? AND method = ?",
		r.URL.Path, r.Method).Scan(&handlerFunc)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Get and execute the handler function
	var sqlCode string
	err = db.QueryRow("SELECT sql_code FROM functions WHERE name = ?", handlerFunc).Scan(&sqlCode)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Execute the SQL function
	var content string
	var mimeType string
	err = db.QueryRow(sqlCode, r.URL.Path).Scan(&content, &mimeType)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
	fmt.Fprint(w, content)
}

func main() {
	reset := flag.Bool("reset", false, "Reset database and reload schema")
	port := flag.String("port", "8080", "Port to listen on")
	dbPath := flag.String("db", "sqlserve.db", "Path to SQLite database file")
	flag.Parse()

	if err := initDB(*dbPath, *reset); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: http.HandlerFunc(handleRequest),
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("Server starting on port %s...", *port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error: %v", err)

	case sig := <-shutdown:
		log.Printf("Got signal: %v", sig)
		log.Println("Shutting down server...")

		// Give outstanding requests 5 seconds to complete.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
	}
}
