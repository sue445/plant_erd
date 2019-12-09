package mysql

type informationSchemaTables struct {
	TableName string `db:"table_name"`
}

type foreignKey struct {
	ToTable    string `db:"to_table"`
	PrimaryKey string `db:"primary_key"`
	Column     string `db:"column"`
	Name       string `db:"name"`
}
