package lib

import (
	"github.com/cockroachdb/errors"
	"github.com/sue445/plant_erd/adapter"
	"github.com/sue445/plant_erd/db"
)

// LoadSchema load schema from adapter
func LoadSchema(adapter adapter.Adapter) (*db.Schema, error) {
	tableNames, err := adapter.GetAllTableNames()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var tables []*db.Table
	for _, tableName := range tableNames {
		table, err := adapter.GetTable(tableName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tables = append(tables, table)
	}

	return db.NewSchema(tables), nil
}
