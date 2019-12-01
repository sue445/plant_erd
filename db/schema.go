package db

import (
	"fmt"
	"strings"
)

// Schema represents schema info
type Schema struct {
	Tables []*Table
}

// ToErd returns ERD formatted schema
func (s *Schema) ToErd() string {
	var lines []string

	for _, table := range s.Tables {
		lines = append(lines, table.ToErd())
	}

	for _, table := range s.Tables {
		for _, foreignKey := range table.ForeignKeys {
			lines = append(lines, fmt.Sprintf("%s }-- %s", table.Name, foreignKey.ToTable))
		}
	}

	return strings.Join(lines, "\n\n")
}
