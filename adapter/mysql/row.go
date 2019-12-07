package mysql

type informationSchemaTables struct {
	TableName string `db:"TABLE_NAME"`
}

type foreignKey struct {
	ToTable    string `db:"to_table"`
	PrimaryKey string `db:"primary_key"`
	Column     string `db:"column"`
	Name       string `db:"name"`
	OnUpdate   string `db:"on_update"`
	OnDelete   string `db:"on_delete"`
}
