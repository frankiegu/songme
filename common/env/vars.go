package env

import "github.com/emredir/songme/datastore"

// Vars wraps env variables like datastore etc.
type Vars struct {
	DB datastore.Datastore
}
