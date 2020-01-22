package oracle

import (
	"database/sql"
	"fmt"
)

type allTables struct {
	TableName string `db:"TABLE_NAME"`
}

type allTabColumns struct {
	ColumnName    string        `db:"COLUMN_NAME"`
	DataType      string        `db:"DATA_TYPE"`
	DataLength    int           `db:"DATA_LENGTH"`
	DataPrecision sql.NullInt32 `db:"DATA_PRECISION"`
	DataScale     sql.NullInt32 `db:"DATA_SCALE"`
	Nullable      string        `db:"NULLABLE"`
}

func (c *allTabColumns) FormatColumnType() string {
	switch c.DataType {
	case "VARCHAR2":
		return fmt.Sprintf("%s(%d)", c.DataType, c.DataLength)
	case "NUMBER":
		if c.DataPrecision.Valid {
			return fmt.Sprintf("%s(%d)", c.DataType, c.DataPrecision.Int32)
		}
	case "FLOAT":
		if c.DataPrecision.Valid {
			if c.DataScale.Valid {
				return fmt.Sprintf("%s(%d,%d)", c.DataType, c.DataPrecision.Int32, c.DataScale.Int32)
			}
			return fmt.Sprintf("%s(%d)", c.DataType, c.DataPrecision.Int32)
		}
	}

	return c.DataType
}

type primaryKeys struct {
	ColumnName string `db:"COLUMN_NAME"`
}

type foreignKey struct {
	ToTable          string `db:"TO_TABLE"`
	ReferencesColumn string `db:"REFERENCES_COLUMN"`
	ColumnName       string `db:"COLUMN_NAME"`
}

type allIndexes struct {
	IndexName  string `db:"INDEX_NAME"`
	Uniqueness string `db:"UNIQUENESS"`
}

type allIndColumns struct {
	ColumnName string `db:"COLUMN_NAME"`
}
