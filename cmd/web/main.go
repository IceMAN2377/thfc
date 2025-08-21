package main

// TODO: implement "NOT FOUND" logic in GetByTITLE
// TODO: implement "ALREADY EXISTS" in PostText

import (
	"fmt"

	"github.com/IceMAN2377/thfc/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"

	v1http "github.com/IceMAN2377/thfc/internal/transport/http"
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

	repo := postgres.NewRepo(db)
	mux := http.NewServeMux()
	v1http.RegEndpoints(logger, mux, repo)

	logger.Info("Starting server")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
