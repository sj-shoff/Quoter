package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"log/slog"
	"quoter/internal/models"
	"quoter/internal/service"
)

var (
	ErrInvalidID    = errors.New("invalid ID format")
	ErrInvalidInput = errors.New("invalid input data")
)

type QuoteHandler struct {
	service service.QuoteServiceInterface
	logger  *slog.Logger
}

func NewQuoteHandler(service service.QuoteServiceInterface, logger *slog.Logger) *QuoteHandler {
	return &QuoteHandler{
		service: service,
		logger:  logger,
	}
}

func (h *QuoteHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /quotes", h.AddQuote)
	mux.HandleFunc("GET /quotes", h.GetAllQuotes)
	mux.HandleFunc("GET /quotes/random", h.GetRandomQuote)
	mux.HandleFunc("DELETE /quotes/{id}", h.DeleteQuote)
}

func (h *QuoteHandler) AddQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Author string `json:"author"`
		Text   string `json:"quote"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "Failed to decode request body", "error", err)
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Author == "" || req.Text == "" {
		h.logger.ErrorContext(ctx, "Missing required fields")
		h.sendError(w, http.StatusBadRequest, "Author and quote text are required")
		return
	}

	quote := models.Quote{
		Author: req.Author,
		Text:   req.Text,
	}

	newQuote, err := h.service.AddQuote(ctx, quote)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to add quote", "error", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to add quote")
		return
	}

	h.sendJSON(w, http.StatusCreated, newQuote)
}

func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	author := r.URL.Query().Get("author")
	var quotes []models.Quote
	var err error

	if author != "" {
		quotes, err = h.service.GetQuotesByAuthor(ctx, author)
	} else {
		quotes, err = h.service.GetAllQuotes(ctx)
	}

	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get quotes", "error", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to get quotes")
		return
	}

	if len(quotes) == 0 {
		h.sendJSON(w, http.StatusOK, []models.Quote{})
		return
	}

	h.sendJSON(w, http.StatusOK, quotes)
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	quote, err := h.service.GetRandomQuote(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get random quote", "error", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to get random quote")
		return
	}

	if quote == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.sendJSON(w, http.StatusOK, quote)
}

func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.logger.ErrorContext(ctx, "Invalid ID format", "error", err)
		h.sendError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if id <= 0 {
		h.logger.WarnContext(ctx, "Invalid ID value", "id", id)
		h.sendError(w, http.StatusBadRequest, "ID must be positive integer")
		return
	}

	if err := h.service.DeleteQuote(ctx, id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.logger.WarnContext(ctx, "Quote not found", "id", id)
			h.sendError(w, http.StatusNotFound, "Quote not found")
			return
		}

		h.logger.ErrorContext(ctx, "Failed to delete quote", "id", id, "error", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to delete quote")
		return
	}

	h.logger.InfoContext(ctx, "Quote deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *QuoteHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.ErrorContext(context.Background(), "Failed to encode response", "error", err)
	}
}

func (h *QuoteHandler) sendError(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error   string `json:"error"`
		Message string `json:"message,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.ErrorContext(context.Background(), "Failed to encode error response", "error", err)
	}
}
