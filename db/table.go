package db

import (
	"fmt"
	"strings"
)

// Table represents table info
type Table struct {
	Name        string
	Columns     []*Column
	ForeignKeys []*ForeignKey
	Indexes     []*Index
}

// ToErd returns ERD formatted table
func (t *Table) ToErd() string {
	lines := []string{
		fmt.Sprintf("entity %s {", t.Name),
	}

	pkColumns := t.GetPrimaryKeyColumns()
	nonPkColumns := t.GetNonPrimaryKeyColumns()

	var area []string

	if len(pkColumns) > 0 {
		var parts []string
		for _, column := range pkColumns {
			parts = append(parts, "  "+column.ToErd())
		}
		area = append(area, strings.Join(parts, "\n"))
	}

	if len(nonPkColumns) > 0 {
		var parts []string
		for _, column := range nonPkColumns {
			parts = append(parts, "  "+column.ToErd())
		}
		area = append(area, strings.Join(parts, "\n"))
	}

	if len(t.Indexes) > 0 {
		var parts []string
		for _, index := range t.Indexes {
			parts = append(parts, "  "+index.ToErd())
		}
		area = append(area, strings.Join(parts, "\n"))
	}

	lines = append(lines, strings.Join(area, "\n  --\n"))

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// GetPrimaryKeyColumns returns Primary key columns
func (t *Table) GetPrimaryKeyColumns() []*Column {
	var columns []*Column
	for _, column := range t.Columns {
		if column.PrimaryKey {
			columns = append(columns, column)
		}
	}
	return columns
}

// GetNonPrimaryKeyColumns returns Non-Primary key columns
func (t *Table) GetNonPrimaryKeyColumns() []*Column {
	var columns []*Column
	for _, column := range t.Columns {
		if !column.PrimaryKey {
			columns = append(columns, column)
		}
	}
	return columns
}
