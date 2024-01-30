package util

import (
	"database/sql"
	"fmt"

	"github.com/Karanth1r3/url-short-learn/internal/config"

	_ "github.com/lib/pq"
)

func ConnectDB(cfg config.DB) (*sql.DB, error) {
	// connect to db with name "postgres"
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Name, cfg.Username, cfg.Password,
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %w", err)
	}
	// set sessions limits
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(0)
	// try to ping db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is not available: %w", err)
	}

	// if everything's ok
	return db, nil
}
