package main

import (
	"encoding/json"
	"fmt"
	"github.com/IceMAN2377/thfc/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var psql *sqlx.DB

func main() {
	var err error

	connStr := "postgres://postgres:secret@localhost:5432/thfc?sslmode=disable"

	psql, err = sqlx.Open("postgres", connStr)
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
	mux.HandleFunc("GET /texts/{title}", GetByTitle)

	log.Print("Starting")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOMEPAGE"))
}

func postText(w http.ResponseWriter, r *http.Request) {
	log.Printf("Database connection status: %v", psql)
	var record models.Record

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		log.Printf("some error: %v", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if record.Title == "" || record.Content == "" {
		http.Error(w, "Empty input", 400)
		return
	}

	stmt, err := psql.Preparex(`INSERT INTO records (title, content) VALUES ($1, $2)`)
	if err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "Internal server error", 500)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.Title, record.Content)
	if err != nil {
		log.Printf("some error: %v", err)
		http.Error(w, "error", 500)
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func GetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	stmt, err := psql.Preparex(`SELECT title, content FROM records WHERE title=$1`)
	if err != nil {
		log.Printf("Error preparing stmt:%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var record models.Record

	if err := stmt.Get(&record, title); err != nil {
		log.Printf("error retrieving the record:%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(record)

}
