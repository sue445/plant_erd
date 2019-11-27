package sqlite3

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // for sql
)

// Provider represents sqlite3 provider
type Provider struct {
	db *sql.DB
}

// Close represents function for close database
type Close func() error

// NewProvider returns a new Provider instance
func NewProvider(name string) (*Provider, Close, error) {
	db, err := sql.Open("sqlite3", name)

	if err != nil {
		return nil, nil, err
	}

	return &Provider{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (p *Provider) GetAllTableNames() ([]string, error) {
	rows, err := p.db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")

	if err != nil {
		return []string{}, nil
	}

	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return []string{}, err
		}
		tables = append(tables, name)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return tables, nil
}
