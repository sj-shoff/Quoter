package service

import (
	"context"
	"log/slog"
	"quoter/internal/models"
	"quoter/internal/storage"
)

type QuoteService struct {
	storage storage.QuoteStorage
	logger  *slog.Logger
}

func NewQuoteService(storage storage.QuoteStorage, logger *slog.Logger) *QuoteService {
	return &QuoteService{
		storage: storage,
		logger:  logger,
	}
}

func (s *QuoteService) AddQuote(ctx context.Context, quote models.Quote) (models.Quote, error) {
	s.logger.DebugContext(ctx, "Adding new quote", "author", quote.Author)
	return s.storage.AddQuote(ctx, quote)
}

func (s *QuoteService) GetAllQuotes(ctx context.Context) ([]models.Quote, error) {
	s.logger.DebugContext(ctx, "Getting all quotes")
	return s.storage.GetAllQuotes(ctx)
}

func (s *QuoteService) GetRandomQuote(ctx context.Context) (*models.Quote, error) {
	s.logger.DebugContext(ctx, "Getting random quote")
	return s.storage.GetRandomQuote(ctx)
}

func (s *QuoteService) GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error) {
	s.logger.DebugContext(ctx, "Getting quotes by author", "author", author)
	return s.storage.GetQuotesByAuthor(ctx, author)
}

func (s *QuoteService) DeleteQuote(ctx context.Context, id int) error {
	s.logger.DebugContext(ctx, "Deleting quote", "id", id)
	return s.storage.DeleteQuote(ctx, id)
}
