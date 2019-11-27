package sqlite3

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // for sql
	"github.com/sue445/plant_erd/provider"
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
		return []string{}, err
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

type pragmaTableInfoRow struct {
	CID        int
	Name       string
	Type       string
	NotNull    bool
	DfltValue  *string
	PrimaryKey bool
}

// GetTable returns table info
func (p *Provider) GetTable(tableName string) (*provider.Table, error) {
	rows, err := p.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	table := provider.Table{
		Name: tableName,
	}

	for rows.Next() {
		var tableInfoRow pragmaTableInfoRow
		err = rows.Scan(&tableInfoRow.CID, &tableInfoRow.Name, &tableInfoRow.Type, &tableInfoRow.NotNull, &tableInfoRow.DfltValue, &tableInfoRow.PrimaryKey)
		if err != nil {
			return nil, err
		}

		column := provider.Column{
			Name:       tableInfoRow.Name,
			Type:       tableInfoRow.Type,
			NotNull:    tableInfoRow.NotNull,
			PrimaryKey: tableInfoRow.PrimaryKey,
		}

		table.Columns = append(table.Columns, column)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	foreignKeys, err := p.getForeignKeys(tableName)
	if err != nil {
		return nil, err
	}

	table.ForeignKeys = foreignKeys

	return &table, nil
}

type pragmaForeignKeyListRow struct {
	ID       int
	Seq      int
	Table    string
	From     string
	To       string
	OnUpdate string
	OnDelete string
	Match    string
}

func (p *Provider) getForeignKeys(tableName string) ([]provider.ForeignKey, error) {
	rows, err := p.db.Query(fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var foreignKeys []provider.ForeignKey

	for rows.Next() {
		var foreignKeyListRow pragmaForeignKeyListRow
		err = rows.Scan(&foreignKeyListRow.ID, &foreignKeyListRow.Seq, &foreignKeyListRow.Table, &foreignKeyListRow.From, &foreignKeyListRow.To,
			&foreignKeyListRow.OnUpdate, &foreignKeyListRow.OnDelete, &foreignKeyListRow.Match)

		if err != nil {
			return nil, err
		}

		foreignKey := provider.ForeignKey{
			Sequence:   foreignKeyListRow.Seq,
			FromColumn: foreignKeyListRow.From,
			ToColumn:   foreignKeyListRow.To,
			ToTable:    foreignKeyListRow.Table,
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return foreignKeys, nil
}
