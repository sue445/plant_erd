package oracle

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8" // for sql
)

// Adapter represents Oracle adapter
type Adapter struct {
	db *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(config *Config) (*Adapter, Close, error) {
	db, err := sqlx.Connect("oci8", config.FormatDSN())

	if err != nil {
		return nil, nil, err
	}

	return &Adapter{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (a *Adapter) GetAllTableNames() ([]string, error) {
	var rows []allTables
	// c.f. https://github.com/rsim/oracle-enhanced/blob/v6.0.0/lib/active_record/connection_adapters/oracle_enhanced/schema_statements.rb#L15
	err := a.db.Select(&rows, `
		SELECT DECODE(table_name, UPPER(table_name), LOWER(table_name), table_name) AS table_name
		FROM all_tables
		WHERE owner = SYS_CONTEXT('userenv', 'current_schema')
		AND secondary = 'N'
		minus
		SELECT DECODE(mview_name, UPPER(mview_name), LOWER(mview_name), mview_name) AS table_name
		FROM all_mviews
		WHERE owner = SYS_CONTEXT('userenv', 'current_schema')
	`)

	if err != nil {
		return []string{}, err
	}

	var tables []string
	for _, row := range rows {
		tables = append(tables, row.TableName)
	}

	return tables, nil
}
