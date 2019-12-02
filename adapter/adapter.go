package adapter

import "github.com/sue445/plant_erd/db"

// Adapter represents database adapter
type Adapter interface {
	GetAllTableNames() ([]string, error)
	GetTable(tableName string) (*db.Table, error)
}
