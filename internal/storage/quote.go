package storage

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"github.com/zonder12120/brandscout-quotebook/internal/model"
)

var ErrNotFound = fmt.Errorf("not found")

type QuoteStorage interface {
	CreateQuote(q *model.Quote) (*model.Quote, error)
	GetQuotesList() ([]*model.Quote, error)
	GetRandomQuote() (*model.Quote, error)
	GetQuotesByAuthor(author string) ([]*model.Quote, error)
	DeleteByID(id int) error
}

type MemoryStorage struct {
	limit  int
	mu     sync.RWMutex
	quotes map[int]*model.Quote
	nextID int
	minID  int
}

func NewInMemory(limitQuotes int) *MemoryStorage {
	return &MemoryStorage{
		limit:  limitQuotes,
		quotes: make(map[int]*model.Quote),
		nextID: 1,
		minID:  1,
	}
}

func (r *MemoryStorage) CreateQuote(q *model.Quote) (*model.Quote, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.quotes) >= r.limit {
		delete(r.quotes, r.minID)
		r.minID++
	}

	q.ID = r.nextID
	r.quotes[q.ID] = q
	r.nextID++

	return q, nil
}

func (r *MemoryStorage) GetQuotesList() ([]*model.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quotes := make([]*model.Quote, 0, len(r.quotes))
	for _, q := range r.quotes {
		quotes = append(quotes, q)
	}

	return quotes, nil
}

func (r *MemoryStorage) GetRandomQuote() (*model.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.quotes) == 0 {
		return nil, ErrNotFound
	}

	ids := make([]int, 0, len(r.quotes))
	for id := range r.quotes {
		ids = append(ids, id)
	}

	randomID := ids[rand.Intn(len(ids))]
	return r.quotes[randomID], nil
}

func (r *MemoryStorage) GetQuotesByAuthor(author string) ([]*model.Quote, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	quotes := make([]*model.Quote, 0, len(r.quotes))

	for _, q := range r.quotes {
		if strings.EqualFold(q.Author, author) {
			quotes = append(quotes, q)
		}
	}
	return quotes, nil
}

func (r *MemoryStorage) DeleteByID(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.quotes[id]; !ok {
		return ErrNotFound
	}

	delete(r.quotes, id)
	return nil
}
