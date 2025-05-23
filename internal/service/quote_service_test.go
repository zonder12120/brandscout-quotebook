package service

import (
	"errors"
	"testing"

	"github.com/zonder12120/brandscout-quotebook/internal/model"
	"github.com/zonder12120/brandscout-quotebook/internal/storage"
)

var (
	errStorage = errors.New("storage error")
	errDB      = errors.New("db error")
)

type mockStorage struct {
	createErr      error
	deleteErr      error
	getRandomErr   error
	getByAuthorErr error
	listErr        error

	createdQuote *model.Quote
	quotesList   []*model.Quote
	authorArg    string
	calledWith   int
}

func (m *mockStorage) CreateQuote(q *model.Quote) (*model.Quote, error) {
	return m.createdQuote, m.createErr
}

func (m *mockStorage) GetQuotesList() ([]*model.Quote, error) {
	return m.quotesList, m.listErr
}

func (m *mockStorage) GetRandomQuote() (*model.Quote, error) {
	return m.createdQuote, m.getRandomErr
}

func (m *mockStorage) GetQuotesByAuthor(author string) ([]*model.Quote, error) {
	m.authorArg = author
	return m.quotesList, m.getByAuthorErr
}

func (m *mockStorage) DeleteByID(id int) error {
	m.calledWith = id
	return m.deleteErr
}

func TestQuoteService(t *testing.T) {
	testQuote := &model.Quote{ID: 1, Author: "Test", Quote: "Test"}
	testQuotes := []*model.Quote{testQuote}

	t.Run("Create", func(t *testing.T) {
		tt := []struct {
			name        string
			mock        *mockStorage
			input       *model.Quote
			expectedErr error
			wantQuote   *model.Quote
		}{
			{
				name:      "success",
				mock:      &mockStorage{createdQuote: testQuote},
				input:     testQuote,
				wantQuote: testQuote,
			},
			{
				name:        "storage error",
				mock:        &mockStorage{createErr: errStorage},
				input:       testQuote,
				expectedErr: errStorage,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				service := NewQuoteService(tc.mock)
				result, err := service.Create(tc.input)

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}

				if tc.wantQuote != nil && result != tc.wantQuote {
					t.Errorf("expected quote %v, got %v", tc.wantQuote, result)
				}
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		tt := []struct {
			name          string
			mock          *mockStorage
			inputID       int
			expectedErr   error
			expectedCalls int
		}{
			{
				name:          "success",
				mock:          &mockStorage{},
				inputID:       1,
				expectedCalls: 1,
			},
			{
				name:          "not found",
				mock:          &mockStorage{deleteErr: storage.ErrNotFound},
				inputID:       2,
				expectedErr:   storage.ErrNotFound,
				expectedCalls: 2,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				service := NewQuoteService(tc.mock)
				err := service.Delete(tc.inputID)

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}

				if tc.mock.calledWith != tc.inputID {
					t.Errorf("expected ID %d, got %d", tc.inputID, tc.mock.calledWith)
				}
			})
		}
	})

	t.Run("List", func(t *testing.T) {
		tt := []struct {
			name        string
			mock        *mockStorage
			expected    []*model.Quote
			expectedErr error
		}{
			{
				name:     "success",
				mock:     &mockStorage{quotesList: testQuotes},
				expected: testQuotes,
			},
			{
				name:        "storage error",
				mock:        &mockStorage{listErr: errDB},
				expectedErr: errDB,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				service := NewQuoteService(tc.mock)
				result, err := service.List()

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}

				if tc.expected != nil && len(result) != len(tc.expected) {
					t.Errorf("expected %d quotes, got %d", len(tc.expected), len(result))
				}
			})
		}
	})

	t.Run("GetRandom", func(t *testing.T) {
		tt := []struct {
			name        string
			mock        *mockStorage
			expected    *model.Quote
			expectedErr error
		}{
			{
				name:     "success",
				mock:     &mockStorage{createdQuote: testQuote},
				expected: testQuote,
			},
			{
				name:        "not found",
				mock:        &mockStorage{getRandomErr: storage.ErrNotFound},
				expectedErr: storage.ErrNotFound,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				service := NewQuoteService(tc.mock)
				result, err := service.GetRandom()

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}

				if tc.expected != nil && result != tc.expected {
					t.Errorf("expected quote %v, got %v", tc.expected, result)
				}
			})
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		tt := []struct {
			name        string
			mock        *mockStorage
			inputAuthor string
			expected    []*model.Quote
			expectedErr error
			expectedArg string
		}{
			{
				name:        "success",
				mock:        &mockStorage{quotesList: testQuotes},
				inputAuthor: "Test",
				expected:    testQuotes,
				expectedArg: "Test",
			},
			{
				name:        "storage error",
				mock:        &mockStorage{getByAuthorErr: errDB},
				inputAuthor: "Error",
				expectedErr: errDB,
				expectedArg: "Error",
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				service := NewQuoteService(tc.mock)
				result, err := service.GetByAuthor(tc.inputAuthor)

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}

				if tc.expected != nil && len(result) != len(tc.expected) {
					t.Errorf("expected %d quotes, got %d", len(tc.expected), len(result))
				}

				if tc.mock.authorArg != tc.expectedArg {
					t.Errorf("expected author %s, got %s", tc.expectedArg, tc.mock.authorArg)
				}
			})
		}
	})
}
