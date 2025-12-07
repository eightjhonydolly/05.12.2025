package test

import (
	"context"
	"testing"
	"time"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/repository"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/service"
)

func TestFullWorkflow_ServiceToRepository(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	svc := service.NewLinkService(repo)

	urls := []string{"httpbin.org", "invalid.test"}
	batch, err := svc.CheckLinks(context.Background(), urls)
	if err != nil {
		t.Fatalf("CheckLinks failed: %v", err)
	}

	if batch.ID != 1 {
		t.Errorf("Expected batch ID 1, got %d", batch.ID)
	}

	if len(batch.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(batch.Links))
	}

	retrieved, err := repo.GetBatch(batch.ID)
	if err != nil {
		t.Fatalf("GetBatch failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Batch not found in repository")
	}

	pdfData, err := svc.GenerateReport([]int{batch.ID})
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	if len(pdfData) == 0 {
		t.Error("Expected PDF data, got empty slice")
	}
}

func TestConcurrentLinkChecking(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	svc := service.NewLinkService(repo)

	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func(id int) {
			urls := []string{"httpbin.org"}
			_, err := svc.CheckLinks(context.Background(), urls)
			if err != nil {
				t.Errorf("Concurrent CheckLinks %d failed: %v", id, err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 3; i++ {
		select {
		case <-done:
		case <-time.After(30 * time.Second):
			t.Fatal("Timeout waiting for concurrent operations")
		}
	}

	batches, err := repo.GetBatches([]int{1, 2, 3})
	if err != nil {
		t.Fatalf("GetBatches failed: %v", err)
	}

	if len(batches) != 3 {
		t.Errorf("Expected 3 batches, got %d", len(batches))
	}
}
