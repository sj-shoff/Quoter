package service

import (
	"context"
	"errors"
	"quoter/internal/models"
)

var ErrNotFound = errors.New("not found")

type QuoteServiceInterface interface {
	AddQuote(ctx context.Context, quote models.Quote) (models.Quote, error)
	GetAllQuotes(ctx context.Context) ([]models.Quote, error)
	GetRandomQuote(ctx context.Context) (*models.Quote, error)
	GetQuotesByAuthor(ctx context.Context, author string) ([]models.Quote, error)
	DeleteQuote(ctx context.Context, id int) error
}
