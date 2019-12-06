package mysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sue445/plant_erd/db"
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
	var rows []informationSchemaTables
	err := a.db.Select(&rows, "SELECT table_name FROM information_schema.tables WHERE table_schema=database() ORDER BY table_name")

	if err != nil {
		return []string{}, err
	}

	var tables []string
	for _, row := range rows {
		tables = append(tables, row.TableName)
	}

	return tables, nil
}

// GetTable returns table info
func (a *Adapter) GetTable(tableName string) (*db.Table, error) {
	table := db.Table{
		Name: tableName,
	}

	var rows []column
	err := a.db.Select(&rows, fmt.Sprintf("SHOW COLUMNS FROM %s", tableName))

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		column := &db.Column{
			Name:       row.Field,
			Type:       row.Type,
			NotNull:    row.Null == "NO",
			PrimaryKey: row.Key == "PRI",
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
	var rows []foreignKey

	// c.f. https://github.com/rails/rails/blob/v6.0.1/activerecord/lib/active_record/connection_adapters/abstract_mysql_adapter.rb#L385-L400
	sql := `SELECT fk.referenced_table_name AS 'to_table',
                   fk.referenced_column_name AS 'primary_key',
                   fk.column_name AS 'column',
                   fk.constraint_name AS 'name',
                   rc.update_rule AS 'on_update',
                   rc.delete_rule AS 'on_delete'
            FROM information_schema.referential_constraints rc
            JOIN information_schema.key_column_usage fk
            USING (constraint_schema, constraint_name)
            WHERE fk.referenced_column_name IS NOT NULL
              AND fk.table_schema = database()
              AND fk.table_name = ?
              AND rc.constraint_schema = database()
              AND rc.table_name = ?`

	err := a.db.Select(&rows, sql, tableName, tableName)
	if err != nil {
		return nil, err
	}

	var foreignKeys []*db.ForeignKey
	for _, row := range rows {
		foreignKey := &db.ForeignKey{
			FromColumn: row.Column,
			ToColumn:   row.PrimaryKey,
			ToTable:    row.ToTable,
		}

		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

func (a *Adapter) getIndexes(tableName string) ([]*db.Index, error) {
	var rows []index
	err := a.db.Select(&rows, fmt.Sprintf("SHOW INDEX FROM %s WHERE Key_name != 'PRIMARY'", tableName))

	if err != nil {
		return nil, err
	}

	var indexes []*db.Index

	currentIndex := ""
	for _, row := range rows {
		if row.KeyName != currentIndex {
			index := db.Index{
				Name:   row.KeyName,
				Unique: row.NonUnique == 0,
			}
			indexes = append(indexes, &index)
			currentIndex = row.KeyName
		}

		last := len(indexes) - 1
		indexes[last].Columns = append(indexes[last].Columns, row.ColumnName)
	}

	return indexes, nil
}
