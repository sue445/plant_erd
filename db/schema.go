package db

import (
	"fmt"
	"github.com/deckarep/golang-set/v2"
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
func (s *Schema) ToErd(showIndex bool) string {
	var lines []string
	tableNames := mapset.NewSet[string]()

	for _, table := range s.Tables {
		lines = append(lines, table.ToErd(showIndex))
		tableNames.Add(table.Name)
	}

	for _, table := range s.Tables {
		for _, foreignKey := range table.ForeignKeys {
			toTable := strings.ToLower(foreignKey.ToTable)
			if tableNames.Contains(toTable) {
				lines = append(lines, fmt.Sprintf("%s }-- %s", table.Name, toTable))
			}
		}
	}

	return strings.Join(lines, "\n\n")
}

// ToMermaid returns Mermaid formatted table
func (s *Schema) ToMermaid(showComment bool) string {
	var lines []string
	tableNames := mapset.NewSet[string]()

	lines = append(lines, "erDiagram")

	for _, table := range s.Tables {
		lines = append(lines, table.ToMermaid(showComment))
		tableNames.Add(table.Name)
	}

	for _, table := range s.Tables {
		for _, foreignKey := range table.ForeignKeys {
			toTable := strings.ToLower(foreignKey.ToTable)
			if tableNames.Contains(toTable) {
				lines = append(lines, fmt.Sprintf("%s ||--o{ %s : owns", toTable, table.Name))
			}
		}
	}

	return strings.Join(lines, "\n\n")
}

// Subset returns subset of a schema
func (s *Schema) Subset(tableName string, distance int) *Schema {
	explorer := NewSchemaExplorer(s)
	tableNames := explorer.Explore(tableName, distance)

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
