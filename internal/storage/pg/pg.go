package pg

import (
	"database/sql"
)

const tablename = "url"

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	storage := &Storage{DB: db}
	return storage
}
