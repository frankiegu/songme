package songme

import (
	"os"
)

var config Config

// Config represents the configuration variables.
type Config struct {
	Port        string
	SongmeAdmin string
	DatabaseURL string
}

// GetConfig returns a copy of the config object since we want config object to be immutable..
func GetConfig() Config {
	return config
}

func init() {
	if os.Getenv("ENV") == "PRODUCTION" && os.Getenv("DATABASE_URL") == "" {
		panic("environment variable DATABASE_URL is not set")
	}

	config = Config{
		Port:        ":" + env("PORT", "8080"),
		SongmeAdmin: env("SONGME_ADMIN", ""),
		DatabaseURL: env("DATABASE_URL", "postgres://username:password@localhost:5432/dbname?sslmode=disable"),
	}

}

func env(key, defvalue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defvalue
}
