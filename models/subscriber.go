package models

import (
	"time"
)

// Subscriber represents our newsletter subscribers
type Subscriber struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}
