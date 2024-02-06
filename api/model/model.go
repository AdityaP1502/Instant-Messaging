package model

import (
	"io"
)

type Model interface {
	FromJSON(r io.Reader, checkRequired bool) error
	ToJSON(checkRequired bool) ([]byte, error)

	// Migrate db functionality to querynator
	// Insert(db *sql.DB) error
	// Update(db *sql.DB, condition map[string]string) error
	// Delete(db *sql.DB) error
	// IsExists(db *sql.DB) (bool, error)

	// Query(db *sql.DB) error
	// JoinQuery(db *sql.DB) error
}
