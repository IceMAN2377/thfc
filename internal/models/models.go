package models

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Record struct {
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type RecordHandler struct {
	db *sqlx.DB
}

func NewRecordHandler(db *sqlx.DB) *RecordHandler {
	return &RecordHandler{
		db: db,
	}
}

func (h *RecordHandler) PostText(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := h.db.Preparex(`INSERT INTO records (title, content) VALUES ($1, $2)`)
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

func (h *RecordHandler) GetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	stmt, err := h.db.Preparex(`SELECT title, content FROM records WHERE title=$1`)
	if err != nil {
		log.Printf("Error preparing stmt:%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var record Record

	if err := stmt.Get(&record, title); err != nil {
		log.Printf("error retrieving the record:%v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(record)

}
