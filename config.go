package songme

import (
	"os"
	"strconv"
)

// Environment variables for Songme.
const (
	EnvPlatform              = "ENV"  // Optional for development
	EnvPort                  = "PORT" // Optional for development
	EnvDatabaseURL           = "DATABASE_URL"
	EnvSongmeAdmin           = "SONGME_ADMIN"
	EnvCookieName            = "COOKIE_NAME"             // Optional
	EnvCookieHashKey         = "COOKIE_HASH_KEY"         // Optional for development
	EnvCookieBlockKey        = "COOKIE_BLOCK_KEY"        // Optional for development
	EnvUsernameLength        = "USERNAME_LENGTH"         // Optional
	EnvPasswordLength        = "PASSWORD_LENGTH"         // Optional
	EnvHashCost              = "HASH_COST"               // Optional
	EnvSongDescriptionLength = "SONG_DESCRIPTION_LENGTH" // Optional
	EnvSongsPerPage          = "SONGS_PER_PAGE"          // Optional
	EnvUsersPerPage          = "USERS_PER_PAGE"          // Optional
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
	Env                   string
	Port                  string
	DatabaseURL           string
	SongmeAdmin           string
	CookieName            string
	CookieHashKey         string
	CookieBlockKey        string
	UsernameLength        int
	PasswordLength        int
	HashCost              int
	SongDescriptionLength int
	SongsPerPage          int
	UsersPerPage          int
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

		SongDescriptionLength: EnvInt(EnvSongDescriptionLength, 280),
		SongsPerPage:          EnvInt(EnvSongsPerPage, 20),

		UsersPerPage: EnvInt(EnvUsersPerPage, 20),
	}
}
