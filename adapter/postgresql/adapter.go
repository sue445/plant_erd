package postgresql

import (
	"github.com/deckarep/golang-set"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sql
	"github.com/sue445/plant_erd/db"
)

// Adapter represents PostgreSQL adapter
type Adapter struct {
	db     *sqlx.DB
	dbName string
}

// Close represents function for close database
type Close func() error

// NewAdapter returns a new Adapter instance
func NewAdapter(config *Config) (*Adapter, Close, error) {
	db, err := sqlx.Connect("postgres", config.FormatDSN())

	if err != nil {
		return nil, nil, err
	}

	return &Adapter{db: db, dbName: config.DBName}, db.Close, nil
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

// GetTable returns table info
func (a *Adapter) GetTable(tableName string) (*db.Table, error) {
	table := db.Table{
		Name: tableName,
	}

	primaryKeyColumns, err := a.getPrimaryKeyColumns(tableName)
	if err != nil {
		return nil, err
	}

	var rows []informationSchemaColumns
	err = a.db.Select(&rows, `
		SELECT column_name,
		       data_type,
		       is_nullable
		FROM information_schema.columns
		WHERE table_catalog = $1 AND table_name = $2
		ORDER BY ordinal_position
	`, a.dbName, tableName)

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		column := &db.Column{
			Name:       row.ColumnName,
			Type:       row.DataType,
			NotNull:    row.IsNullable == "NO",
			PrimaryKey: primaryKeyColumns.Contains(row.ColumnName),
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

func (a *Adapter) getPrimaryKeyColumns(tableName string) (mapset.Set, error) {
	var rows []primaryKeys

	err := a.db.Select(&rows, `
		SELECT ccu.column_name as COLUMN_NAME
		FROM information_schema.table_constraints tc,
		     information_schema.constraint_column_usage ccu
		WHERE tc.table_catalog=$1
		AND tc.table_name=$2
		AND tc.constraint_type='PRIMARY KEY'
		AND tc.table_catalog=ccu.table_catalog
		AND tc.table_schema=ccu.table_schema
		AND tc.table_name=ccu.table_name
		AND tc.constraint_name=ccu.constraint_name
	`, a.dbName, tableName)

	if err != nil {
		return nil, err
	}

	columns := mapset.NewSet()
	for _, row := range rows {
		columns.Add(row.ColumnName)
	}

	return columns, nil
}

func (a *Adapter) getForeignKeys(tableName string) ([]*db.ForeignKey, error) {
	var rows []foreignKey

	// c.f. https://github.com/rails/rails/blob/v6.0.1/activerecord/lib/active_record/connection_adapters/postgresql/schema_statements.rb#L483
	err := a.db.Select(&rows, `
		SELECT t2.oid::regclass::text AS to_table, a1.attname AS column, a2.attname AS primary_key, c.conname AS name
		FROM pg_constraint c
		JOIN pg_class t1 ON c.conrelid = t1.oid
		JOIN pg_class t2 ON c.confrelid = t2.oid
		JOIN pg_attribute a1 ON a1.attnum = c.conkey[1] AND a1.attrelid = t1.oid
		JOIN pg_attribute a2 ON a2.attnum = c.confkey[1] AND a2.attrelid = t2.oid
		JOIN pg_namespace t3 ON c.connamespace = t3.oid
		WHERE c.contype = 'f'
		  AND t1.relname = $1
		  AND t3.nspname = $2
		ORDER BY c.conname
	`, tableName, "public")

	if err != nil {
		return nil, err
	}

	var foreignKeys []*db.ForeignKey
	for _, row := range rows {
		foreignKey := &db.ForeignKey{
			FromColumn: row.Column,
			ToTable:    row.ToTable,
			ToColumn:   row.PrimaryKey,
		}
		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

func (a *Adapter) getIndexes(tableName string) ([]*db.Index, error) {
	// c.f. https://github.com/rails/rails/blob/v6.0.1/activerecord/lib/active_record/connection_adapters/postgresql/schema_statements.rb#L89
	var rows []indexes
	err := a.db.Select(&rows, `
		SELECT distinct i.relname, d.indisunique, d.indkey, t.oid
		FROM pg_class t
		INNER JOIN pg_index d ON t.oid = d.indrelid
		INNER JOIN pg_class i ON d.indexrelid = i.oid
		LEFT JOIN pg_namespace n ON n.oid = i.relnamespace
		WHERE i.relkind = 'i'
		  AND d.indisprimary = 'f'
		  AND t.relname = $1
		ORDER BY i.relname
	`, tableName)

	if err != nil {
		return nil, err
	}

	var indexes []*db.Index
	for _, row := range rows {
		columns, err := a.getIndexColumns(row.Oid, row.Indkeys())
		if err != nil {
			return nil, err
		}

		index := &db.Index{
			Name:    row.Relname,
			Unique:  row.Indisunique,
			Columns: columns,
		}
		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (a *Adapter) getIndexColumns(oid int, indkeys []int) ([]string, error) {
	// c.f. https://github.com/rails/rails/blob/v6.0.1/activerecord/lib/active_record/connection_adapters/postgresql/schema_statements.rb#L119
	sql := "SELECT a.attnum AS attnum, a.attname AS attname FROM pg_attribute a WHERE a.attrelid = ? AND a.attnum IN (?)"

	query, args, err := sqlx.In(sql, oid, indkeys)
	if err != nil {
		return nil, err
	}

	query = a.db.Rebind(query)

	var rows []pgAttribute
	err = a.db.Select(&rows, query, args...)

	if err != nil {
		return nil, err
	}

	columnNames := map[int]string{}
	for _, row := range rows {
		columnNames[row.Attnum] = row.Attname
	}

	var columns []string
	for _, indkey := range indkeys {
		columns = append(columns, columnNames[int(indkey)])
	}

	return columns, nil
}
