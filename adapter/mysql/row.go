package mysql

import "database/sql"

type informationSchemaTables struct {
	TableName string `db:"TABLE_NAME"`
}

type column struct {
	Field   string         `db:"Field"`
	Type    string         `db:"Type"`
	Null    string         `db:"Null"`
	Key     string         `db:"Key"`
	Default sql.NullString `db:"Default"`
	Extra   string         `db:"Extra"`
}

type foreignKey struct {
	ToTable    string `db:"to_table"`
	PrimaryKey string `db:"primary_key"`
	Column     string `db:"column"`
	Name       string `db:"name"`
	OnUpdate   string `db:"on_update"`
	OnDelete   string `db:"on_delete"`
}

type index struct {
	Table        string         `db:"Table"`
	NonUnique    int            `db:"Non_unique"`
	KeyName      string         `db:"Key_name"`
	SeqInIndex   int            `db:"Seq_in_index"`
	ColumnName   string         `db:"Column_name"`
	Collation    string         `db:"Collation"`
	Cardinality  string         `db:"Cardinality"`
	SubPart      sql.NullInt32  `db:"Sub_part"`
	Packed       sql.NullString `db:"Packed"`
	Null         string         `db:"Null"`
	IndexType    string         `db:"Index_type"`
	Comment      string         `db:"Comment"`
	IndexComment string         `db:"Index_comment"`
	Visible      string         `db:"Visible"`
	Expression   sql.NullString `db:"Expression"`
}
