package context

import (
	"context"
	"log"

	"github.com/emredir/songme/models"
)

type ctxKey string

const (
	userKey ctxKey = "user"
)

// WithUser returns a new context with the given user assigned to it.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User will return the user assigned to the context.
func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	if val == nil {
		return nil
	}
	user, ok := val.(*models.User)
	if !ok {
		log.Println("[context.User]: user value set incorrectly. val=", val)
		return nil
	}
	return user
}
