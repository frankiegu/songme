package models

import (
	"time"
)

// Subscriber represents our newsletter subscribers.
type Subscriber struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// NewSubscriber returns new subscriber.
func NewSubscriber(name, email string) *Subscriber {
	return &Subscriber{
		Name:  name,
		Email: email,
	}
}
