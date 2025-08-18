package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Record struct {
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

var psql *sql.DB

func main() {
	var err error

	connStr := "postgres://postgres:secret@localhost:5432/thfc?sslmode=disable"

	psql, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to prepare DB")
	}

	err = psql.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}

	fmt.Println("Success to DB")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("POST /texts/", postText)

	log.Print("Starting")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOMEPAGE"))
}

func postText(w http.ResponseWriter, r *http.Request) {
	log.Printf("Database connection status: %v", psql)
	var record Record

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		log.Printf("some error: %v", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if record.Title == "" || record.Content == "" {
		http.Error(w, "Empty input", 400)
		return
	}

	stmt, err := psql.Prepare(`INSERT INTO records (title, content) VALUES ($1, $2)`)
	if err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.Title, record.Content)
	if err != nil {
		log.Printf("some error: %v", err)
		http.Error(w, "DALBAEB", 500)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}
