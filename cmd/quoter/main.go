package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quoter/internal/http-server/handler"
	middleware "quoter/internal/http-server/middleware"
	"quoter/internal/service"
	"quoter/internal/storage/inmemory"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	storage, err := inmemory.NewInMemoryStorage()
	if err != nil {
		log.Error("Failed to init storage")
		os.Exit(1)
	}

	quoteService := service.NewQuoteService(storage, log)

	quoteHandler := handler.NewQuoteHandler(quoteService, log)

	mux := http.NewServeMux()
	quoteHandler.RegisterRoutes(mux)

	handlerWithLogger := middleware.New(log)(mux)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithLogger,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Failed to start server")
		}
	}()

	log.Info("Server started")

	<-done
	log.Info("Stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Failed to stop server")

		return
	}

	log.Info("Server stopped")
}
