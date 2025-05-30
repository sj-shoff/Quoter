package tests

import (
	"context"
	"quoter/internal/models"
	"quoter/internal/service"
	"quoter/internal/storage"
	"testing"
)

type mockStorage struct {
	quotes []models.Quote
}

func (m *mockStorage) AddQuote(ctx context.Context, quote models.Quote) (models.Quote, error) {
	quote.ID = len(m.quotes) + 1
	m.quotes = append(m.quotes, quote)
	return quote, nil
}

func (m *mockStorage) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	return m.quotes, nil
}

func (m *mockStorage) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	if len(m.quotes) == 0 {
		return nil, storage.ErrNotFound
	}
	return &m.quotes[0], nil
}

func (m *mockStorage) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	var result []models.Quote
	for _, q := range m.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result, nil
}

func (m *mockStorage) DeleteQuote(ctx context.Context, id int) error {
	for i, q := range m.quotes {
		if q.ID == id {
			m.quotes = append(m.quotes[:i], m.quotes[i+1:]...)
			return nil
		}
	}
	return storage.ErrNotFound
}

func TestQuoteService(t *testing.T) {
	ctx := context.Background()
	storage := &mockStorage{}
	service := service.NewQuoteService(storage, nil)

	t.Run("Add Quote", func(t *testing.T) {
		_, err := service.AddQuote(ctx, models.Quote{Author: "Author", Text: "Quote"})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get All Quotes", func(t *testing.T) {
		_, err := service.GetAllQuotes(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Delete Quote", func(t *testing.T) {
		quote, _ := service.AddQuote(ctx, models.Quote{Author: "ToDelete", Text: "Delete me"})
		err := service.DeleteQuote(ctx, quote.ID)
		if err != nil {
			t.Fatal(err)
		}
	})
}
