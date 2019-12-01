package sqlite3

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // for sql
	"github.com/sue445/plant_erd/provider"
)

// Provider represents sqlite3 provider
type Provider struct {
	db *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewProvider returns a new Provider instance
func NewProvider(name string) (*Provider, Close, error) {
	db, err := sqlx.Connect("sqlite3", name)

	if err != nil {
		return nil, nil, err
	}

	return &Provider{db: db}, db.Close, nil
}

// GetAllTableNames returns all table names in database
func (p *Provider) GetAllTableNames() ([]string, error) {
	var rows []sqliteMaster
	err := p.db.Select(&rows, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")

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
func (p *Provider) GetTable(tableName string) (*provider.Table, error) {
	table := provider.Table{
		Name: tableName,
	}

	var rows []tableInfo
	err := p.db.Select(&rows, fmt.Sprintf("PRAGMA table_info(%s)", tableName))

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		column := &provider.Column{
			Name:       row.Name,
			Type:       row.Type,
			NotNull:    row.NotNull,
			PrimaryKey: row.PrimaryKey,
		}

		table.Columns = append(table.Columns, column)
	}

	foreignKeys, err := p.getForeignKeys(tableName)
	if err != nil {
		return nil, err
	}

	table.ForeignKeys = foreignKeys

	return &table, nil
}

func (p *Provider) getForeignKeys(tableName string) ([]*provider.ForeignKey, error) {
	var rows []foreignKeyList
	err := p.db.Select(&rows, fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName))

	if err != nil {
		return nil, err
	}

	var foreignKeys []*provider.ForeignKey
	for _, row := range rows {
		foreignKey := &provider.ForeignKey{
			Sequence:   row.Seq,
			FromColumn: row.From,
			ToColumn:   row.To,
			ToTable:    row.Table,
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}
