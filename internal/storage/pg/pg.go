package pg

import (
	"database/sql"
	"fmt"

	"github.com/Karanth1r3/url-short-learn/internal/storage"
	"github.com/lib/pq"
)

const tablename = "shorts"

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) *Storage {
	storage := &Storage{DB: db}
	return storage
}

// Tries to save url with generated alias to storage
func (s *Storage) SaveURL(urlToSave, alias string) error {
	const erDesc = "storage.pg.SaveURL"
	stmt, err := s.DB.Prepare("INSERT INTO shorts(url, alias) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", erDesc, err)
	}
	//TODO - check constraints (uniqueness)
	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			fmt.Println("Severity:", err.Severity)
			fmt.Println("Code:", err.Code)
			fmt.Println("Message:", err.Message)
			fmt.Println("Detail:", err.Detail)
			fmt.Println("Hint:", err.Hint)
			fmt.Println("Position:", err.Position)
			fmt.Println("InternalPosition:", err.InternalPosition)
			fmt.Println("Where:", err.Where)
			fmt.Println("Schema:", err.Schema)
			fmt.Println("Table:", err.Table)
			fmt.Println("Column:", err.Column)
			fmt.Println("DataTypeName:", err.DataTypeName)
			fmt.Println("Constraint:", err.Constraint)
			fmt.Println("File:", err.File)
			fmt.Println("Line:", err.Line)
			fmt.Println("Routine:", err.Routine)
		}
		// db already has that alias/url
		return storage.ErrURLExists
	}
	// Get id of the last inserted record. Not supported everywhere
	/*
		id, err := res.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("%s: %w", erDesc, err)
		}
	*/
	return nil
}
