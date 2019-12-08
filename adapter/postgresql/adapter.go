package postgresql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sql
)

// Adapter represents PostgreSQL adapter
type Adapter struct {
	db *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(config *Config) (*Adapter, Close, error) {
	db, err := sqlx.Connect("postgres", config.FormatDSN())

	if err != nil {
		return nil, nil, err
	}

	return &Adapter{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (a *Adapter) GetAllTableNames() ([]string, error) {
	var rows []pgStatUserTables
	err := a.db.Select(&rows, "SELECT relname FROM pg_stat_user_tables ORDER BY relname")

	if err != nil {
		return []string{}, err
	}

	var tables []string
	for _, row := range rows {
		tables = append(tables, row.Relname)
	}

	return tables, nil
}
