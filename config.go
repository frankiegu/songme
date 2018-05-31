package songme

import (
	"os"
	"strconv"
)

// Environment variables for Songme.
const (
	EnvPlatform       = "ENV"
	EnvPort           = "PORT"
	EnvDatabaseURL    = "DATABASE_URL"
	EnvSongmeAdmin    = "SONGME_ADMIN"
	EnvCookieName     = "COOKIE_NAME"
	EnvCookieHashKey  = "COOKIE_HASH_KEY"
	EnvCookieBlockKey = "COOKIE_BLOCK_KEY"
	EnvUsernameLength = "USERNAME_LENGTH"
	EnvPasswordLength = "PASSWORD_LENGTH"
	EnvHashCost       = "HASH_COST"
)

// Env retrieves the value of the environment variable named by the key.
func Env(key, defvalue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defvalue
}

// EnvInt retrieves the value of the environment variable named by the key.
func EnvInt(key string, defvalue int) int {
	value, _ := strconv.Atoi(os.Getenv(key))
	if value != 0 {
		return value
	}
	return defvalue
}

/*
	App Configuration
*/

var config Config

// Config represents the configuration variables.
type Config struct {
	Env            string
	Port           string
	DatabaseURL    string
	SongmeAdmin    string
	CookieName     string
	CookieHashKey  string
	CookieBlockKey string
	UsernameLength int
	PasswordLength int
	HashCost       int
}

// GetConfig returns a copy of the config object since we want config object to be immutable.
func GetConfig() Config {
	return config
}

func init() {
	if os.Getenv(EnvPlatform) == "PRODUCTION" && os.Getenv(EnvDatabaseURL) == "" {
		panic("environment variable DATABASE_URL is not set")
	}

	config = Config{
		Env:         Env(EnvPlatform, "DEVELOPMENT"),
		Port:        ":" + Env(EnvPort, "8080"),
		DatabaseURL: Env(EnvDatabaseURL, "postgres://username:password@localhost:5432/dbname?sslmode=disable"),

		SongmeAdmin: Env(EnvSongmeAdmin, ""),

		CookieName:     Env(EnvCookieName, "songmeSession"),
		CookieHashKey:  Env(EnvCookieHashKey, ""),
		CookieBlockKey: Env(EnvCookieBlockKey, ""),

		UsernameLength: EnvInt(EnvUsernameLength, 5),
		PasswordLength: EnvInt(EnvPasswordLength, 6),
		HashCost:       EnvInt(EnvHashCost, 10),
	}
}
