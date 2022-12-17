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
func (t *Table) ToErd(showIndex bool) string {
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

	if showIndex && len(t.Indexes) > 0 {
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

// ToMermaid returns Mermaid formatted table
func (t *Table) ToMermaid(showComment bool) string {
	lines := []string{
		fmt.Sprintf("%s {", t.Name),
	}

	for _, column := range t.Columns {
		var parts []string
		parts = append(parts, column.ToMermaid())

		if showComment {
			key := t.mermaidColumnKey(column)
			if key != "" {
				parts = append(parts, key)
			}

			comment := t.mermaidColumnComment(column)
			if comment != "" {
				parts = append(parts, comment)
			}
		}

		line := "  " + strings.Join(parts, " ")
		lines = append(lines, line)
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func (t *Table) mermaidColumnKey(column *Column) string {
	if column.PrimaryKey {
		return "PK"
	}

	for _, foreignKey := range t.ForeignKeys {
		if foreignKey.FromColumn == column.Name {
			return "FK"
		}
	}

	return ""
}

func (t *Table) mermaidColumnComment(column *Column) string {
	parts := []string{}
	if column.NotNull {
		parts = append(parts, "not null")
	}

	if len(parts) == 0 {
		return ""
	}
	return fmt.Sprintf("\"%s\"", strings.Join(parts, " "))
}
