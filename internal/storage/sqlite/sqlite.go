package sqlite

import (
	"database/sql"

	// Driver is used by sql package, not directly by code => _ is necessary
	_ "github.com/lib/p"
)

type Storage struct {
	DB *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// String for additional info to put into errors
	//const op = "storage.sqlite.New"
	// Opening connection to db
	/*
		db, err := sql.Open("sqlite3", storagePath)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fmt.Println("before")
		// Query for db creation
		stmt := (`CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias on url(alias);`)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		_, err = db.Exec(stmt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	*/
	//return &Storage{db}, nil
	return nil, nil
}
