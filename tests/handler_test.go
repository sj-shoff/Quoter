package tests

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"quoter/internal/http-server/handler"
	"quoter/internal/models"
	"quoter/internal/service"
	"quoter/internal/storage/inmemory"
	"strconv"
	"testing"
)

func setupHandler() *handler.QuoteHandler {
	storage, _ := inmemory.NewInMemoryStorage()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	service := service.NewQuoteService(storage, logger)
	return handler.NewQuoteHandler(service, logger)
}

func addTestQuote(h *handler.QuoteHandler, author, text string) *models.Quote {
	body := bytes.NewBufferString(`{"author": "` + author + `", "quote": "` + text + `"}`)
	req := httptest.NewRequest("POST", "/quotes", body)
	w := httptest.NewRecorder()
	h.AddQuote(w, req)

	var quote models.Quote
	json.NewDecoder(w.Body).Decode(&quote)
	return &quote
}

func TestAddQuoteHandler(t *testing.T) {
	h := setupHandler()

	t.Run("Valid Request", func(t *testing.T) {
		body := bytes.NewBufferString(`{"author": "Test", "quote": "Hello"}`)
		req := httptest.NewRequest("POST", "/quotes", body)
		w := httptest.NewRecorder()

		h.AddQuote(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var quote models.Quote
		json.NewDecoder(w.Body).Decode(&quote)
		if quote.Text != "Hello" {
			t.Error("Incorrect quote saved")
		}
	})

	t.Run("Invalid Request", func(t *testing.T) {
		body := bytes.NewBufferString(`{"author": ""}`)
		req := httptest.NewRequest("POST", "/quotes", body)
		w := httptest.NewRecorder()

		h.AddQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestGetQuotesHandler(t *testing.T) {
	h := setupHandler()

	addTestQuote(h, "Author1", "Quote1")
	addTestQuote(h, "Author2", "Quote2")

	t.Run("Get All Quotes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes", nil)
		w := httptest.NewRecorder()

		h.GetAllQuotes(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var quotes []models.Quote
		json.NewDecoder(w.Body).Decode(&quotes)
		if len(quotes) != 2 {
			t.Errorf("Expected 2 quotes, got %d", len(quotes))
		}
	})

	t.Run("Filter by Author", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes?author=Author1", nil)
		w := httptest.NewRecorder()

		h.GetAllQuotes(w, req)

		var quotes []models.Quote
		json.NewDecoder(w.Body).Decode(&quotes)
		if len(quotes) != 1 || quotes[0].Author != "Author1" {
			t.Errorf("Filter by author failed: %v", quotes)
		}
	})
}

func TestRandomQuoteHandler(t *testing.T) {
	h := setupHandler()

	t.Run("No Quotes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		h.GetRandomQuote(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})

	t.Run("With Quotes", func(t *testing.T) {
		addTestQuote(h, "R", "Random")

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		h.GetRandomQuote(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestDeleteQuoteHandler(t *testing.T) {
	h := setupHandler()

	quote := addTestQuote(h, "ToDelete", "Delete me")

	t.Run("Valid Delete", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(quote.ID), nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/invalid", nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/quotes/999", nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
