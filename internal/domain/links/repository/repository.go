package repository

import (
	"sync"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
)

type LinkRepository interface {
	SaveBatch(batch *model.LinkBatch) error
	GetBatch(id int) (*model.LinkBatch, error)
	GetBatches(ids []int) ([]*model.LinkBatch, error)
	GetNextID() int
}

type InMemoryLinkRepository struct {
	mu      sync.RWMutex
	batches map[int]*model.LinkBatch
	nextID  int
}

func NewInMemoryLinkRepository() *InMemoryLinkRepository {
	return &InMemoryLinkRepository{
		batches: make(map[int]*model.LinkBatch),
		nextID:  1,
	}
}

func (r *InMemoryLinkRepository) SaveBatch(batch *model.LinkBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.batches[batch.ID] = batch
	return nil
}

func (r *InMemoryLinkRepository) GetBatch(id int) (*model.LinkBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	batch, exists := r.batches[id]
	if !exists {
		return nil, nil
	}
	return batch, nil
}

func (r *InMemoryLinkRepository) GetBatches(ids []int) ([]*model.LinkBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var batches []*model.LinkBatch
	for _, id := range ids {
		if batch, exists := r.batches[id]; exists {
			batches = append(batches, batch)
		}
	}
	return batches, nil
}

func (r *InMemoryLinkRepository) GetNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++
	return id
}
