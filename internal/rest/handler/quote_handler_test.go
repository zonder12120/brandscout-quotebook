package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/zonder12120/brandscout-quotebook/internal/model"
	"github.com/zonder12120/brandscout-quotebook/internal/storage"
	"github.com/zonder12120/brandscout-quotebook/pkg/logger"
)

type mockService struct {
	createErr      error
	deleteErr      error
	getRandomErr   error
	getByAuthorErr error
	createdQuote   *model.Quote
	quotesList     []*model.Quote
}

func (m *mockService) Create(q *model.Quote) (*model.Quote, error) {
	return m.createdQuote, m.createErr
}

func (m *mockService) List() ([]*model.Quote, error) {
	return m.quotesList, nil
}

func (m *mockService) GetRandom() (*model.Quote, error) {
	return m.createdQuote, m.getRandomErr
}

func (m *mockService) GetByAuthor(author string) ([]*model.Quote, error) {
	return m.quotesList, m.getByAuthorErr
}

func (m *mockService) Delete(id int) error {
	return m.deleteErr
}

func TestHandler(t *testing.T) {
	log := logger.New("debug")

	t.Run("Create validation failure", func(t *testing.T) {
		reqBody := []byte(`{"author": "", "quote": ""}`)
		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h := New(&mockService{}, log)
		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var resp map[string]string
		err := json.NewDecoder(rec.Body).Decode(&resp)
		if err != nil {
			t.Errorf("expected success decode, got %v", err)
		}
		if resp["error"] != errEmptyAuthorOrQuote {
			t.Errorf("unexpected error message: %s", resp["error"])
		}
	})

	t.Run("Create success", func(t *testing.T) {
		expectedQuote := &model.Quote{ID: 1, Author: "Test", Quote: "Test"}
		reqBody := []byte(`{"author": "Test", "quote": "Test"}`)
		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h := New(&mockService{createdQuote: expectedQuote}, log)
		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}

		var quote model.Quote
		err := json.NewDecoder(rec.Body).Decode(&quote)
		if err != nil {
			t.Errorf("expected success decode, got %v", err)
		}
		if quote.ID != expectedQuote.ID {
			t.Errorf("expected ID %d, got %d", expectedQuote.ID, quote.ID)
		}
	})

	t.Run("Delete invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rec := httptest.NewRecorder()

		h := New(&mockService{}, log)
		h.Delete(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("Delete not found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		h := New(&mockService{deleteErr: storage.ErrNotFound}, log)
		h.Delete(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("Get random quote not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes/random", nil)
		rec := httptest.NewRecorder()

		h := New(&mockService{getRandomErr: storage.ErrNotFound}, log)
		h.Random(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("Filter by author success", func(t *testing.T) {
		expectedQuotes := []*model.Quote{
			{ID: 1, Author: "Test", Quote: "Test"},
		}
		req := httptest.NewRequest("GET", "/quotes?author=Test", nil)
		rec := httptest.NewRecorder()

		h := New(&mockService{quotesList: expectedQuotes}, log)
		h.FilterByAuthor(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var quotes []model.Quote
		err := json.NewDecoder(rec.Body).Decode(&quotes)
		if err != nil {
			t.Errorf("expected success decode, got %v", err)
		}
		if len(quotes) != 1 {
			t.Errorf("expected 1 quote, got %d", len(quotes))
		}
	})
}
