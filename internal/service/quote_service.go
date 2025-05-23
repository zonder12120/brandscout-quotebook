package service

import (
	"github.com/zonder12120/brandscout-quotebook/internal/model"
	"github.com/zonder12120/brandscout-quotebook/internal/storage"
)

type Quote interface {
	Create(q *model.Quote) (*model.Quote, error)
	List() ([]*model.Quote, error)
	GetRandom() (*model.Quote, error)
	GetByAuthor(author string) ([]*model.Quote, error)
	Delete(id int) error
}

type QuoteService struct {
	store storage.QuoteStorage
}

func NewQuoteService(store storage.QuoteStorage) *QuoteService {
	return &QuoteService{store: store}
}

func (s *QuoteService) Create(q *model.Quote) (*model.Quote, error) {
	return s.store.CreateQuote(q)
}

func (s *QuoteService) List() ([]*model.Quote, error) {
	return s.store.GetQuotesList()
}

func (s *QuoteService) GetRandom() (*model.Quote, error) {
	return s.store.GetRandomQuote()
}

func (s *QuoteService) GetByAuthor(author string) ([]*model.Quote, error) {
	return s.store.GetQuotesByAuthor(author)
}

func (s *QuoteService) Delete(id int) error {
	return s.store.DeleteByID(id)
}
