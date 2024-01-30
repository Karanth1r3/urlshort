package sqlite

import (
	"database/sql"
	"fmt"

	// Driver is used by sql package, not directly by code => _ is necessary
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// String for additional info to put into errors
	const op = "storage.sqlite.New"
	// Opening connection to db
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// Query for db creation
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE
		url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias on url(alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db}, nil
}
