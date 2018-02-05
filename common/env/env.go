package env

import "github.com/emre-demir/songme/datastore"

// Env wraps env variables like datastore etc.
type Env struct {
	DB datastore.Datastore
}
