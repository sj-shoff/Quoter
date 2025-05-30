package inmemory

import (
	"context"
	"math/rand"
	"quoter/internal/models"
	"quoter/internal/storage"
	"sync"
	"sync/atomic"
	"time"
)

type inMemoryStorage struct {
	quotes []models.Quote
	mu     sync.RWMutex
	nextID atomic.Int64
	rand   *rand.Rand
}

func NewInMemoryStorage() storage.QuoteStorage {
	src := rand.NewSource(time.Now().UnixNano())
	return &inMemoryStorage{
		quotes: make([]models.Quote, 0),
		rand:   rand.New(src),
	}
}

func (s *inMemoryStorage) AddQuote(ctx context.Context, quote models.Quote) (models.Quote, error) {
	quote.ID = int(s.nextID.Add(1))

	s.mu.Lock()
	defer s.mu.Unlock()

	s.quotes = append(s.quotes, quote)
	return quote, nil
}

func (s *inMemoryStorage) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Quote, len(s.quotes))
	copy(result, s.quotes)
	return result, nil
}

func (s *inMemoryStorage) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.quotes) == 0 {
		return nil, nil
	}
	idx := s.rand.Intn(len(s.quotes))
	return &s.quotes[idx], nil
}

func (s *inMemoryStorage) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Quote, 0)
	for _, q := range s.quotes {
		if q.Author == author {
			result = append(result, q)
		}
	}
	return result, nil
}

func (s *inMemoryStorage) DeleteQuote(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, q := range s.quotes {
		if q.ID == id {
			s.quotes = append(s.quotes[:i], s.quotes[i+1:]...)
			return nil
		}
	}
	return storage.ErrNotFound
}
