package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Adapter represents sqlite3 adapter
type Adapter struct {
	db *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(config *mysql.Config) (*Adapter, Close, error) {
	db, err := sqlx.Connect("mysql", config.FormatDSN())

	if err != nil {
		return nil, nil, err
	}

	return &Adapter{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (a *Adapter) GetAllTableNames() ([]string, error) {
	var tables []string
	return tables, nil
}
