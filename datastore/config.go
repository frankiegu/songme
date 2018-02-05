package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config holds configuration information belongs to database.
type Config struct {
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
	SSLMode string
}

// DefaultConfig holds the configuration information taken from 'config.json' file.
var DefaultConfig Config

// PQConn formats and returns dataSourceName for PostgresSQL.
func (config Config) PQConn() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.SSLMode,
	)
}

// EnsurePQReady checks if the PostgreSQL database is ready.
// Executes each statement, if any error occurs then returns error.
func (config Config) EnsurePQReady(statements []string) error {
	db, err := sql.Open("postgres", config.PQConn())
	if err != nil {
		log.Println("[EnsurePQReady]:", err)
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Println("[EnsurePQReady]:", err)
		return err
	}
	for _, statement := range statements {
		_, err = db.Exec(statement)
		if err != nil {
			log.Println("[EnsurePQReady]:", err)
			return err
		}
	}
	return nil
}

func init() {
	file, err := os.Open("datastore/config.json")
	if err != nil {
		log.Fatal("[DATASTORE - INIT]", err)
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&DefaultConfig)
	if err != nil {
		log.Fatal("[DATASTORE - INIT:]", err)
	}
}
