package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/zonder12120/brandscout-quotebook/internal/model"
	"github.com/zonder12120/brandscout-quotebook/internal/service"
	"github.com/zonder12120/brandscout-quotebook/internal/storage"
	"github.com/zonder12120/brandscout-quotebook/pkg/logger"
)

const (
	errInvalidRequestPayload = "invalid request payload"
	errInvalidResponse       = "invalid response"
	errCreateQuote           = "failed to create quote"
	errGetQuotes             = "failed to get quotes"
	errGetRandomQuote        = "failed to get random quote"
	errGetByAuthor           = "failed to get quotes by author"
	errGetID                 = "failed to get quote id"
	errDeleteQuote           = "failed to delete quote"
	errQuoteNotFound         = "quote not found"
	errEmptyAuthorOrQuote    = "author and quote must be non-empty"
	errEmptyAuthor           = "author param required"
)

type QuoteHandler struct {
	service service.Quote
	logger  *logger.Logger
}

func New(service service.Quote, logger *logger.Logger) *QuoteHandler {
	return &QuoteHandler{
		service: service,
		logger:  logger,
	}
}

func (h *QuoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var quote model.Quote
	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		h.respondError(w, http.StatusBadRequest, errInvalidRequestPayload, err)
		return
	}

	if strings.TrimSpace(quote.Author) == "" || strings.TrimSpace(quote.Quote) == "" {
		h.respondError(w, http.StatusBadRequest, errEmptyAuthorOrQuote, nil)
		return
	}

	created, err := h.service.Create(&quote)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, errCreateQuote, err)
		return
	}

	respondJSON(w, http.StatusCreated, created)
}

func (h *QuoteHandler) List(w http.ResponseWriter, _ *http.Request) {
	quotes, err := h.service.List()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, errGetQuotes, err)
		return
	}

	respondJSON(w, http.StatusOK, quotes)
}

func (h *QuoteHandler) Random(w http.ResponseWriter, _ *http.Request) {
	quote, err := h.service.GetRandom()
	if errors.Is(err, storage.ErrNotFound) {
		h.respondError(w, http.StatusNotFound, errQuoteNotFound, err)
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, errGetRandomQuote, err)
		return
	}

	respondJSON(w, http.StatusOK, quote)
}

func (h *QuoteHandler) FilterByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if strings.TrimSpace(author) == "" {
		h.respondError(w, http.StatusBadRequest, errEmptyAuthor, nil)
		return
	}

	quote, err := h.service.GetByAuthor(author)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, errGetByAuthor, err)
		return
	}

	respondJSON(w, http.StatusOK, quote)
}

func (h *QuoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]

	id, err := strconv.Atoi(strID)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, errGetID, err)
		return
	}

	err = h.service.Delete(id)
	if errors.Is(err, storage.ErrNotFound) {
		h.respondError(w, http.StatusNotFound, errQuoteNotFound, err)
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, errDeleteQuote, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, errInvalidResponse, http.StatusInternalServerError)
	}
}

func (h *QuoteHandler) respondError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		h.logger.Error().Err(err).Msg(message)
	} else {
		h.logger.Warn().Msg(message)
	}

	respondJSON(w, status, map[string]string{"error": message})
}

func (h *QuoteHandler) logDebug(r *http.Request) {
	h.logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.String())
}
