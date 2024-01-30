package pg

import (
	"database/sql"
	"fmt"

	"github.com/Karanth1r3/url-short-learn/internal/config"
)

const tablename = "url"

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	storage := &Storage{DB: db}
	return storage
}

func InitStorage(db *sql.DB, cfg config.DB) error {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", cfg.Host, cfg.Port, tablename, cfg.Username, cfg.Password)
	db, err := sql.Open("postgres", connStr)
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE
		url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias on url(alias);}`)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
