package subscription

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateSubscription(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	sub := Subscription{
		ID:     "sub-1001",
		UserID: "user-500",
		Plan:   "premium",
	}

	err := service.Create(ctx, sub)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateSubscriptionMissingID(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	sub := Subscription{
		UserID: "user-500",
		Plan:   "premium",
	}

	err := service.Create(ctx, sub)

	if err != ErrInvalidSubscription {
		t.Fatalf("expected %v, got %v", ErrInvalidSubscription, err)
	}
}

func TestCreateSubscriptionDuplicateID(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	sub := Subscription{
		ID:     "sub-1001",
		UserID: "user-500",
		Plan:   "premium",
	}

	err := service.Create(ctx, sub)
	if err != nil {
		t.Fatalf("expected first create to succeed, got %v", err)
	}

	err = service.Create(ctx, sub)

	if err != ErrDuplicateID {
		t.Fatalf("expected %v, got %v", ErrDuplicateID, err)
	}
}

func TestGetSubscription(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	sub := Subscription{
		ID:     "sub-1001",
		UserID: "user-500",
		Plan:   "premium",
	}

	if err := service.Create(ctx, sub); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, exists := service.Get(ctx, "sub-1001")

	if !exists {
		t.Fatal("expected subscription to exist")
	}

	if got.ID != "sub-1001" {
		t.Fatalf("expected sub-1001, got %s", got.ID)
	}
}

func TestConcurrentCreateAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(i int) {
			sub := Subscription{
				ID:     fmt.Sprintf("sub-%d", i),
				UserID: fmt.Sprintf("user-%d", i),
				Plan:   "premium",
			}

			_ = service.Create(ctx, sub)
			_, _ = service.Get(ctx, sub.ID)

			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestCreateSubscriptionCancelledContext(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sub := Subscription{
		ID:     "sub-2001",
		UserID: "user-900",
		Plan:   "premium",
	}

	err := service.Create(ctx, sub)

	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestCancelSubscription(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	ctx := context.Background()

	sub := Subscription{
		ID:     "sub-3001",
		UserID: "user-700",
		Plan:   "premium",
	}

	if err := service.Create(ctx, sub); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := service.Cancel(ctx, "sub-3001"); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}

	got, exists := service.Get(ctx, "sub-3001")
	if !exists {
		t.Fatal("expected subscription to exist")
	}

	if got.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", got.Status)
	}
}
