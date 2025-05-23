package storage

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/zonder12120/brandscout-quotebook/internal/model"
)

func TestMemoryStorage(t *testing.T) {
	t.Run("CreateQuote and GetQuotesList", func(t *testing.T) {
		s := NewInMemory(10)
		q := &model.Quote{Author: "Test", Quote: "Test quote"}

		created, err := s.CreateQuote(q)
		if err != nil {
			t.Fatalf("createQuote failed: %v", err)
		}

		if created.ID != 1 {
			t.Errorf("expected ID 1, got %d", created.ID)
		}

		list, err := s.GetQuotesList()
		if err != nil {
			t.Fatalf("getQuotesList failed: %v", err)
		}

		if len(list) != 1 {
			t.Errorf("expected 1 quote, got %d", len(list))
		}
	})

	t.Run("Quote rotation when limit exceeded", func(t *testing.T) {
		limit := 3
		s := NewInMemory(limit)

		for i := 0; i < limit+2; i++ {
			_, err := s.CreateQuote(&model.Quote{
				Author: "Author",
				Quote:  "Quote " + strconv.Itoa(i+1),
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		list, _ := s.GetQuotesList()
		if len(list) != limit {
			t.Fatalf("expected %d quotes, got %d", limit, len(list))
		}

		for _, q := range list {
			if q.ID < 3 {
				t.Errorf("unexpected quote with ID %d", q.ID)
			}
		}
	})

	t.Run("GetRandomQuote", func(t *testing.T) {
		s := NewInMemory(10)

		_, err := s.GetRandomQuote()
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}

		quotes := make([]*model.Quote, 5)
		for i := 0; i < 5; i++ {
			q := &model.Quote{Author: "A", Quote: "Q" + strconv.Itoa(i+1)}
			quotes[i], _ = s.CreateQuote(q)
		}

		found := make(map[int]bool)
		for i := 0; i < 100; i++ {
			q, err := s.GetRandomQuote()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			found[q.ID] = true
		}

		if len(found) != 5 {
			t.Errorf("expected all quotes to be returned, got %d unique", len(found))
		}
	})

	t.Run("DeleteByID", func(t *testing.T) {
		s := NewInMemory(10)

		err := s.DeleteByID(999)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}

		created, _ := s.CreateQuote(&model.Quote{Author: "A", Quote: "Q"})
		if err := s.DeleteByID(created.ID); err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		if len(s.quotes) != 0 {
			t.Error("Quote not deleted")
		}
	})

	t.Run("GetQuotesByAuthor", func(t *testing.T) {
		s := NewInMemory(10)

		authors := []string{"AuthorA", "AuthorB", "authorA"}
		for i, author := range authors {
			_, err := s.CreateQuote(&model.Quote{
				Author: author,
				Quote:  "Quote " + strconv.Itoa(i+1),
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		quotes, err := s.GetQuotesByAuthor("authora")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(quotes) != 2 {
			t.Errorf("expected 2 quotes, got %d", len(quotes))
		}

		for _, q := range quotes {
			if !strings.EqualFold(q.Author, "authora") {
				t.Errorf("unexpected author: %s", q.Author)
			}
		}

		quotes, err = s.GetQuotesByAuthor("Unknown")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(quotes) != 0 {
			t.Errorf("expected 0 quotes, got %d", len(quotes))
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		s := NewInMemory(100)
		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				_, err := s.CreateQuote(&model.Quote{
					Author: "Author" + strconv.Itoa(n),
					Quote:  "Quote" + strconv.Itoa(n),
				})
				if err != nil {
					t.Errorf("create error: %v", err)
				}
			}(i)
		}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := s.GetQuotesList()
				if err != nil {
					t.Errorf("list error: %v", err)
				}
			}()
		}

		wg.Wait()
		list, _ := s.GetQuotesList()
		if len(list) != 100 {
			t.Errorf("expected 100 quotes, got %d", len(list))
		}
	})
}
