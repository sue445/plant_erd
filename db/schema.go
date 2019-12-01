package db

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"strings"
)

// Schema represents schema info
type Schema struct {
	Tables []*Table
}

// ToErd returns ERD formatted schema
func (s *Schema) ToErd() string {
	var lines []string
	tableNames := mapset.NewSet()

	for _, table := range s.Tables {
		lines = append(lines, table.ToErd())
		tableNames.Add(table.Name)
	}

	for _, table := range s.Tables {
		for _, foreignKey := range table.ForeignKeys {
			if tableNames.Contains(foreignKey.ToTable) {
				lines = append(lines, fmt.Sprintf("%s }-- %s", table.Name, foreignKey.ToTable))
			}
		}
	}

	return strings.Join(lines, "\n\n")
}

// SurroundingTablesWithin returns surrounding tables from table
func (s *Schema) SurroundingTablesWithin(tableName string, distance int) []string {
	explorer := NewSchemaExplorer(s)
	return explorer.Explore(tableName, distance)
}
