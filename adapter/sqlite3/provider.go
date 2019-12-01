package sqlite3

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // for sql
	"github.com/sue445/plant_erd/db"
)

// Adapter represents sqlite3 adapter
type Adapter struct {
	db *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(name string) (*Adapter, Close, error) {
	db, err := sqlx.Connect("sqlite3", name)

	if err != nil {
		return nil, nil, err
	}

	return &Adapter{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (a *Adapter) GetAllTableNames() ([]string, error) {
	var rows []sqliteMaster
	err := a.db.Select(&rows, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")

	if err != nil {
		return []string{}, err
	}

	var tables []string
	for _, row := range rows {
		tables = append(tables, row.Name)
	}

	return tables, nil
}

// GetTable returns table info
func (a *Adapter) GetTable(tableName string) (*db.Table, error) {
	table := db.Table{
		Name: tableName,
	}

	var rows []tableInfo
	err := a.db.Select(&rows, fmt.Sprintf("PRAGMA table_info(%s)", tableName))

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		column := &db.Column{
			Name:       row.Name,
			Type:       row.Type,
			NotNull:    row.NotNull,
			PrimaryKey: row.PrimaryKey,
		}

		table.Columns = append(table.Columns, column)
	}

	foreignKeys, err := a.getForeignKeys(tableName)
	if err != nil {
		return nil, err
	}

	table.ForeignKeys = foreignKeys

	return &table, nil
}

func (a *Adapter) getForeignKeys(tableName string) ([]*db.ForeignKey, error) {
	var rows []foreignKeyList
	err := a.db.Select(&rows, fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName))

	if err != nil {
		return nil, err
	}

	var foreignKeys []*db.ForeignKey
	for _, row := range rows {
		foreignKey := &db.ForeignKey{
			Sequence:   row.Seq,
			FromColumn: row.From,
			ToColumn:   row.To,
			ToTable:    row.Table,
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}
