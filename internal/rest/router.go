package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zonder12120/brandscout-quotebook/internal/rest/handler"
	"github.com/zonder12120/brandscout-quotebook/internal/rest/middleware"
	"github.com/zonder12120/brandscout-quotebook/pkg/logger"
)

func NewRouter(h *handler.QuoteHandler, logger *logger.Logger) http.Handler {
	r := mux.NewRouter()

	r.Use(middleware.Logging(logger))

	r.HandleFunc("/quotes", h.Create).Methods("POST")
	r.HandleFunc("/quotes", h.FilterByAuthor).Methods("GET").Queries("author", "{author}")
	r.HandleFunc("/quotes", h.List).Methods("GET")
	r.HandleFunc("/quotes/random", h.Random).Methods("GET")
	r.HandleFunc("/quotes/{id:[0-9]+}", h.Delete).Methods("DELETE")

	return r
}
