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

// NewSchema returns a new Schema instance
func NewSchema(tables []*Table) *Schema {
	return &Schema{Tables: tables}
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

// Subset returns subset of a schema
func (s *Schema) Subset(tableName string, distance int) *Schema {
	tableNames := s.SurroundingTablesWithin(tableName, distance)

	var tables []*Table
	for _, tableName := range tableNames {
		table := s.findTable(tableName)

		if table != nil {
			tables = append(tables, table)
		}
	}

	return NewSchema(tables)
}

func (s *Schema) findTable(tableName string) *Table {
	for _, table := range s.Tables {
		if table.Name == tableName {
			return table
		}
	}
	return nil
}
