package tests

import (
	"context"
	"quoter/internal/models"
	"quoter/internal/storage/inmemory"
	"strconv"
	"sync"
	"testing"
)

func TestInMemoryStorage(t *testing.T) {
	ctx := context.Background()
	storage, err := inmemory.NewInMemoryStorage()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	t.Run("Add and Get Quote", func(t *testing.T) {
		quote := models.Quote{Author: "Test Author", Text: "Test Quote"}
		added, err := storage.AddQuote(ctx, quote)
		if err != nil {
			t.Fatal(err)
		}

		if added.ID <= 0 {
			t.Error("Expected positive ID")
		}

		quotes, err := storage.GetAllQuotes(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if len(quotes) != 1 || quotes[0].Text != "Test Quote" {
			t.Error("Failed to retrieve added quote")
		}
	})

	t.Run("Get Random Quote", func(t *testing.T) {
		_, err := storage.GetRandomQuote(ctx)
		if err != nil {
			t.Error("Failed to get random quote")
		}
	})

	t.Run("Delete Quote", func(t *testing.T) {
		quotes, _ := storage.GetAllQuotes(ctx)
		if len(quotes) == 0 {
			t.Fatal("No quotes available")
		}

		err := storage.DeleteQuote(ctx, quotes[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		quotesAfter, _ := storage.GetAllQuotes(ctx)
		if len(quotesAfter) != 0 {
			t.Error("Quote not deleted")
		}
	})
	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				storage.AddQuote(ctx, models.Quote{
					Author: strconv.Itoa(i),
					Text:   "Quote" + strconv.Itoa(i),
				})
			}(i)
		}
		wg.Wait()

		quotes, _ := storage.GetAllQuotes(ctx)
		if len(quotes) != 100 {
			t.Errorf("Expected 100 quotes, got %d", len(quotes))
		}
	})
}
