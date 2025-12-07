package service

import (
	"context"
	"testing"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/repository"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
)

func TestLinkService_CheckLinks(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	service := NewLinkService(repo)

	urls := []string{"google.com", "invalid.test"}
	batch, err := service.CheckLinks(context.Background(), urls)

	if err != nil {
		t.Fatalf("CheckLinks failed: %v", err)
	}

	if batch == nil {
		t.Fatal("Expected batch, got nil")
	}

	if len(batch.Links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(batch.Links))
	}

	if batch.ID != 1 {
		t.Errorf("Expected batch ID 1, got %d", batch.ID)
	}
}

func TestLinkService_GenerateReport(t *testing.T) {
	repo := repository.NewInMemoryLinkRepository()
	service := NewLinkService(repo)

	batch := &model.LinkBatch{
		ID: 1,
		Links: []model.LinkCheck{
			{URL: "google.com", Status: model.StatusAvailable},
		},
	}
	repo.SaveBatch(batch)

	pdfData, err := service.GenerateReport([]int{1})
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	if len(pdfData) == 0 {
		t.Error("Expected PDF data, got empty slice")
	}
}