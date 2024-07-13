package sqlite3

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // for sql
	"github.com/sue445/plant_erd/db"
)

// Adapter represents sqlite3 adapter
type Adapter struct {
	DB *sqlx.DB
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(name string) (*Adapter, Close, error) {
	db, err := sqlx.Connect("sqlite3", name)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &Adapter{DB: db}, db.Close, nil
}

func toBool(i int64) bool {
	return i != 0
}

// GetAllTableNames returns all table names in database
func (a *Adapter) GetAllTableNames() ([]string, error) {
	var rows []sqliteMaster
	err := a.DB.Select(&rows, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")

	if err != nil {
		return []string{}, errors.WithStack(err)
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

	rows, err := a.DB.Queryx(fmt.Sprintf("PRAGMA table_info(%s)", tableName))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}

	table.ForeignKeys = foreignKeys

	indexes, err := a.getIndexes(tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	table.Indexes = indexes

	return &table, nil
}

func (a *Adapter) getForeignKeys(tableName string) ([]*db.ForeignKey, error) {
	rows, err := a.DB.Queryx(fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var foreignKeys []*db.ForeignKey
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		toColumn := ""
		if row["to"] == nil {
			// NOTE: If `to` is NULL, implicitly equals `id` (maybe...)
			// c.f. https://github.com/diesel-rs/diesel/issues/1535
			toColumn = "id"
		} else {
			toColumn = row["to"].(string)
		}

		foreignKey := &db.ForeignKey{
			FromColumn: row["from"].(string),
			ToColumn:   toColumn,
			ToTable:    row["table"].(string),
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

func (a *Adapter) getIndexes(tableName string) ([]*db.Index, error) {
	rows, err := a.DB.Queryx(fmt.Sprintf("PRAGMA index_list(%s)", tableName))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var indexes []*db.Index
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		index := &db.Index{
			Name:   row["name"].(string),
			Unique: row["unique"].(int64) != 0,
		}

		columns, err := a.getIndexColumns(index.Name)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		index.Columns = columns

		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (a *Adapter) getIndexColumns(indexName string) ([]string, error) {
	rows, err := a.DB.Queryx(fmt.Sprintf("PRAGMA index_info(%s)", indexName))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var columns []string
	for rows.Next() {
		row := map[string]interface{}{}
		err := rows.MapScan(row)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		columns = append(columns, row["name"].(string))
	}

	return columns, nil
}
