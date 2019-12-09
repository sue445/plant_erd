package postgresql

import (
	"strconv"
	"strings"
)

type pgStatUserTables struct {
	Relname string `db:"relname"`
}

type informationSchemaColumns struct {
	ColumnName string `db:"column_name"`
	DataType   string `db:"data_type"`
	IsNullable string `db:"is_nullable"`
}

type primaryKeys struct {
	ColumnName string `db:"column_name"`
}

type foreignKey struct {
	ToTable    string `db:"to_table"`
	Column     string `db:"column"`
	PrimaryKey string `db:"primary_key"`
	Name       string `db:"name"`
}

type indexes struct {
	Relname     string `db:"relname"`
	Indisunique bool   `db:"indisunique"`
	Indkey      string `db:"indkey"`
	Oid         int    `db:"oid"`
}

func (i *indexes) Indkeys() []int {
	keys := strings.Split(i.Indkey, " ")

	indkeys := []int{}
	for _, key := range keys {
		i, _ := strconv.Atoi(key)
		indkeys = append(indkeys, i)
	}

	return indkeys
}

type pgAttribute struct {
	Attnum  int    `db:"attnum"`
	Attname string `db:"attname"`
}
