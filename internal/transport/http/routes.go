package http

import (
	"github.com/IceMAN2377/thfc/internal/repository"
	"log/slog"
	"net/http"
)

func RegEndpoints(logger *slog.Logger, mux *http.ServeMux, repo repository.Repository) {

	app := newRecordHandler(repo, logger)

	mux.HandleFunc("POST /texts/", app.PostText)
	mux.HandleFunc("GET /texts/{title}", app.GetByTitle)

}
