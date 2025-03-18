package server

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	// Import the SQLite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// Global variables.
var (
	CloudLinks *Links
	DB         *sql.DB
)

// Represents the JSON links structure.
type Links struct {
	ErrorPage string `json:"error"`
}

// Initialise server port, cloud-links and database (DB).
func Initialise() bool {
	if initialisePort() {
		return false
	}
	initialiseLinks()
	initialiseDB()
	return true
}

// Checks for the "-print-port" flag, for Makefile target
func initialisePort() bool {
	printPort := flag.Bool("print-port", false, "Print a random available port and exit")
	flag.Parse()

	if *printPort {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find a random port: %v\n", err)
			os.Exit(1)
		}
		defer listener.Close()

		fmt.Println(listener.Addr().(*net.TCPAddr).Port)
		return true
	}

	return false
}

// Open the JSON file and unmarshal it's content into CloudLinks struct.
func initialiseLinks() {
	content, err := os.ReadFile("./server/cloudLinks.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &CloudLinks)
	if err != nil {
		return
	}
}

// create/open DB and create tables if they aren't already created.
func initialiseDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Failed to open SQLite database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("connection to the database is dead.", err)
	}

	// Set connection pool configurations
	DB.SetMaxOpenConns(10)                 
	DB.SetMaxIdleConns(5)                  
	DB.SetConnMaxLifetime(5 * time.Minute)

	content, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		log.Fatal("Failed to get database tables:", err)
	}

	if _, err := DB.Exec(string(content)); err != nil {
		log.Fatal("Failed to create database tables:", err)
	}
}
