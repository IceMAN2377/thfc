package http

import (
	"encoding/json"
	"github.com/IceMAN2377/thfc/internal/models"
	"github.com/IceMAN2377/thfc/internal/repository"
	"log"
	"log/slog"
	"net/http"
)

func newRecordHandler(repo repository.Repository, logger *slog.Logger) *recordHandler {
	return &recordHandler{
		repo:   repo,
		logger: logger,
	}
}

type recordHandler struct {
	repo   repository.Repository
	logger *slog.Logger
}

func (h *recordHandler) PostText(w http.ResponseWriter, r *http.Request) {
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

	if err := h.repo.PostText(&record); err != nil {
		h.logger.Error("error writing to DB")
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *recordHandler) GetByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	record, err := h.repo.GetByTitle(title)
	if err != nil {
		h.logger.Error("error retrieving from DB")
		return
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(record)
}
