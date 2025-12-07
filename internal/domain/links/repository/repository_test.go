package repository

import (
	"testing"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
)

func TestInMemoryLinkRepository_SaveAndGetBatch(t *testing.T) {
	repo := NewInMemoryLinkRepository()

	batch := &model.LinkBatch{
		ID: 1,
		Links: []model.LinkCheck{
			{URL: "google.com", Status: model.StatusAvailable},
		},
	}

	err := repo.SaveBatch(batch)
	if err != nil {
		t.Fatalf("SaveBatch failed: %v", err)
	}

	retrieved, err := repo.GetBatch(1)
	if err != nil {
		t.Fatalf("GetBatch failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected batch, got nil")
	}

	if retrieved.ID != 1 {
		t.Errorf("Expected ID 1, got %d", retrieved.ID)
	}
}

func TestInMemoryLinkRepository_GetNextID(t *testing.T) {
	repo := NewInMemoryLinkRepository()

	id1 := repo.GetNextID()
	id2 := repo.GetNextID()

	if id1 != 1 {
		t.Errorf("Expected first ID to be 1, got %d", id1)
	}

	if id2 != 2 {
		t.Errorf("Expected second ID to be 2, got %d", id2)
	}
}

func TestInMemoryLinkRepository_GetBatches(t *testing.T) {
	repo := NewInMemoryLinkRepository()

	batch1 := &model.LinkBatch{ID: 1}
	batch2 := &model.LinkBatch{ID: 2}

	repo.SaveBatch(batch1)
	repo.SaveBatch(batch2)

	batches, err := repo.GetBatches([]int{1, 2})
	if err != nil {
		t.Fatalf("GetBatches failed: %v", err)
	}

	if len(batches) != 2 {
		t.Errorf("Expected 2 batches, got %d", len(batches))
	}
}