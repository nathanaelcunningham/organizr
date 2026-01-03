package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	// Enable WAL mode for better concurrency
	db, err := sql.Open("sqlite3", dataSourceName+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings for better concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
