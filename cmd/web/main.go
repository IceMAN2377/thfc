package main

import (
	"fmt"
	"github.com/IceMAN2377/thfc/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var err error

	connStr := "postgres://postgres:secret@localhost:5432/thfc?sslmode=disable"

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logger.Error("Failed to prepare DB")
	}

	if err = db.Ping(); err != nil {
		logger.Error("Failed to connect DB")
	}

	fmt.Println("Success to DB")

	app := models.NewRecordHandler(db)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /texts/", app.PostText)
	mux.HandleFunc("GET /texts/{title}", app.GetByTitle)

	log.Print("Starting")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
