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

func toBool(i int64) bool {
	return i != 0
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

	rows, err := a.db.Queryx(fmt.Sprintf("PRAGMA table_info(%s)", tableName))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, err
		}

		column := &db.Column{
			Name:       row["name"].(string),
			Type:       row["type"].(string),
			NotNull:    toBool(row["notnull"].(int64)),
			PrimaryKey: toBool(row["pk"].(int64)),
		}

		table.Columns = append(table.Columns, column)
	}

	foreignKeys, err := a.getForeignKeys(tableName)
	if err != nil {
		return nil, err
	}

	table.ForeignKeys = foreignKeys

	indexes, err := a.getIndexes(tableName)
	if err != nil {
		return nil, err
	}

	table.Indexes = indexes

	return &table, nil
}

func (a *Adapter) getForeignKeys(tableName string) ([]*db.ForeignKey, error) {
	rows, err := a.db.Queryx(fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName))

	if err != nil {
		return nil, err
	}

	var foreignKeys []*db.ForeignKey
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, err
		}

		foreignKey := &db.ForeignKey{
			FromColumn: row["from"].(string),
			ToColumn:   row["to"].(string),
			ToTable:    row["table"].(string),
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

func (a *Adapter) getIndexes(tableName string) ([]*db.Index, error) {
	rows, err := a.db.Queryx(fmt.Sprintf("PRAGMA index_list(%s)", tableName))

	if err != nil {
		return nil, err
	}

	var indexes []*db.Index
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, err
		}

		index := &db.Index{
			Name:   row["name"].(string),
			Unique: row["unique"].(int64) != 0,
		}

		columns, err := a.getIndexColumns(index.Name)

		if err != nil {
			return nil, err
		}

		index.Columns = columns

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (a *Adapter) getIndexColumns(indexName string) ([]string, error) {
	rows, err := a.db.Queryx(fmt.Sprintf("PRAGMA index_info(%s)", indexName))

	if err != nil {
		return nil, err
	}

	var columns []string
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, err
		}

		columns = append(columns, row["name"].(string))
	}

	return columns, nil
}
