package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quoter/internal/http-server/handler"
	"quoter/internal/http-server/middleware"
	"quoter/internal/service"
	"quoter/internal/storage/inmemory"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	store := inmemory.NewInMemoryStorage()

	quoteService := service.NewQuoteService(store, logger)

	quoteHandler := handler.NewQuoteHandler(quoteService, logger)

	mux := http.NewServeMux()
	quoteHandler.RegisterRoutes(mux)

	stack := middleware.New(logger)(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: stack,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
		}
	}()
	logger.Info("Server started on :8080")

	<-done
	logger.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}
}
