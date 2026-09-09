package subscription

import (
	"context"
	"errors"
	"sync"
)

type MemoryRepository struct {
	mu            sync.RWMutex
	subscriptions map[string]Subscription
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		subscriptions: make(map[string]Subscription),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, sub Subscription) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.subscriptions[sub.ID]; exists {
		return ErrDuplicateID
	}

	r.subscriptions[sub.ID] = sub
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (Subscription, bool) {
	select {
	case <-ctx.Done():
		return Subscription{}, false
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, exists := r.subscriptions[id]
	return sub, exists
}

func (r *MemoryRepository) Update(ctx context.Context, sub Subscription) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.subscriptions[sub.ID]; !exists {
		return errors.New("subscription not found")
	}

	r.subscriptions[sub.ID] = sub
	return nil
}
