package subscription

import (
	"context"
	"errors"
)

var (
	ErrInvalidSubscription = errors.New("invalid subscription")
	ErrDuplicateID         = errors.New("subscription already exists")
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
	return s.repo.Get(ctx, id)
}
