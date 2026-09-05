package subscription

import "context"

type Repository interface {
	Create(ctx context.Context, sub Subscription) error
	Get(ctx context.Context, id string) (Subscription, bool)
}
