package subscription

import (
	"context"
	"errors"
)

var (
	ErrInvalidSubscription = errors.New("invalid subscription")
	ErrDuplicateID         = errors.New("subscription already exists")
	ErrNotFound            = errors.New("subscription not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, sub Subscription) error {
	if sub.ID == "" || sub.UserID == "" || sub.Plan == "" {
		return ErrInvalidSubscription
	}

	sub.Status = "active"

	return s.repo.Create(ctx, sub)
}

func (s *Service) Get(ctx context.Context, id string) (Subscription, bool) {
	if id == "" {
		return Subscription{}, false
	}

	return s.repo.Get(ctx, id)
}

func (s *Service) Cancel(ctx context.Context, id string) error {
	sub, exists := s.repo.Get(ctx, id)
	if !exists {
		return ErrNotFound
	}

	if sub.Status == "cancelled" {
		return nil
	}

	sub.Status = "cancelled"

	return s.repo.Update(ctx, sub)
}

func (s *Service) Pause(ctx context.Context, id string) error {
	sub, exists := s.repo.Get(ctx, id)
	if !exists {
		return ErrNotFound
	}

	if sub.Status == "cancelled" {
		return errors.New("cancelled subscription cannot be paused")
	}

	if sub.Status == "paused" {
		return nil
	}

	sub.Status = "paused"

	return s.repo.Update(ctx, sub)
}
