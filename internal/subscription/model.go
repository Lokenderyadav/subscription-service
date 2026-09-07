package subscription

import "time"

type Subscription struct {
	ID        string
	UserID    string
	Plan      string
	Status    string
	CreatedAt time.Time
}