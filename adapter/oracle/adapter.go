package oracle

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8" // for sql
	"github.com/sue445/plant_erd/db"
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

// GetTable returns table info
func (a *Adapter) GetTable(tableName string) (*db.Table, error) {
	table := db.Table{
		Name: tableName,
	}

	primaryKeyColumns, err := a.getPrimaryKeyColumns(tableName)
	if err != nil {
		return nil, err
	}

	sql := `
		SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, DATA_PRECISION, DATA_SCALE, NULLABLE
		FROM ALL_TAB_COLUMNS
		WHERE TABLE_NAME = UPPER(?)
		AND owner = SYS_CONTEXT('userenv', 'current_schema')
	`
	stmt, err := a.db.Preparex(a.db.Rebind(sql))
	if err != nil {
		return nil, err
	}

	var rows []allTabColumns
	err = stmt.Select(&rows, tableName)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		column := &db.Column{
			Name:       row.ColumnName,
			Type:       row.FormatColumnType(),
			NotNull:    row.Nullable == "N",
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
	// c.f. https://github.com/rsim/oracle-enhanced/blob/v6.0.0/lib/active_record/connection_adapters/oracle_enhanced_adapter.rb#L612
	sql := `
		SELECT cc.column_name
		FROM all_constraints c, all_cons_columns cc
		WHERE c.owner = SYS_CONTEXT('userenv', 'current_schema')
		AND c.table_name = UPPER(?)
		AND c.constraint_type = 'P'
		AND cc.owner = c.owner
		AND cc.constraint_name = c.constraint_name
		order by cc.position
	`

	stmt, err := a.db.Preparex(a.db.Rebind(sql))
	if err != nil {
		return nil, err
	}

	var rows []primaryKeys
	err = stmt.Select(&rows, tableName)
	defer stmt.Close()

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
	// c.f. https://github.com/rsim/oracle-enhanced/blob/v6.0.0/lib/active_record/connection_adapters/oracle_enhanced/schema_statements.rb#L544
	sql := `
            SELECT r.table_name to_table
                  ,rc.column_name references_column
                  ,cc.column_name
              FROM all_constraints c, all_cons_columns cc,
                   all_constraints r, all_cons_columns rc
             WHERE c.owner = SYS_CONTEXT('userenv', 'current_schema')
               AND c.table_name = UPPER(?)
               AND c.constraint_type = 'R'
               AND cc.owner = c.owner
               AND cc.constraint_name = c.constraint_name
               AND r.constraint_name = c.r_constraint_name
               AND r.owner = c.owner
               AND rc.owner = r.owner
               AND rc.constraint_name = r.constraint_name
               AND rc.position = cc.position
            ORDER BY to_table, column_name, references_column
	`

	stmt, err := a.db.Preparex(a.db.Rebind(sql))
	if err != nil {
		return nil, err
	}

	var rows []foreignKey
	err = stmt.Select(&rows, tableName)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	var foreignKeys []*db.ForeignKey

	for _, row := range rows {
		foreignKey := &db.ForeignKey{
			FromColumn: row.ColumnName,
			ToColumn:   row.ReferencesColumn,
			ToTable:    row.ToTable,
		}
		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

func (a *Adapter) getIndexes(tableName string) ([]*db.Index, error) {
	// c.f. https://github.com/rsim/oracle-enhanced/blob/v6.0.0/lib/active_record/connection_adapters/oracle_enhanced/schema_statements.rb#L91
	sql := `
		SELECT index_name, uniqueness
		FROM all_indexes i
		WHERE owner = SYS_CONTEXT('userenv', 'current_schema') 
		AND table_owner = SYS_CONTEXT('userenv', 'current_schema') 
		AND table_name = UPPER(?)
		AND NOT EXISTS (
			SELECT uc.index_name
			FROM all_constraints uc
			WHERE uc.index_name = i.index_name AND uc.owner = i.owner AND uc.constraint_type = 'P'
		)
		ORDER BY table_name
	`

	stmt, err := a.db.Preparex(a.db.Rebind(sql))
	if err != nil {
		return nil, err
	}

	var rows []allIndexes
	err = stmt.Select(&rows, tableName)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}
	var indexes []*db.Index
	for _, row := range rows {
		columns, err := a.getIndexColumns(row.IndexName)
		if err != nil {
			return nil, err
		}

		index := &db.Index{
			Name:    row.IndexName,
			Unique:  row.Uniqueness == "UNIQUE",
			Columns: columns,
		}
		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (a *Adapter) getIndexColumns(indexName string) ([]string, error) {
	// c.f. https://github.com/rsim/oracle-enhanced/blob/v6.0.0/lib/active_record/connection_adapters/oracle_enhanced/schema_statements.rb#L91
	sql := "SELECT column_name FROM all_ind_columns WHERE index_name = ? ORDER BY column_position"

	stmt, err := a.db.Preparex(a.db.Rebind(sql))
	if err != nil {
		return nil, err
	}

	var rows []allIndColumns
	err = stmt.Select(&rows, indexName)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	var columns []string

	for _, row := range rows {
		columns = append(columns, row.ColumnName)
	}

	return columns, nil
}
