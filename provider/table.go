package provider

import (
	"fmt"
	"strings"
)

// Table represents table info
type Table struct {
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
}

// ToErd returns ERD formatted table
func (t *Table) ToErd() string {
	lines := []string{
		fmt.Sprintf("entity %s {", t.Name),
	}

	pkColumns := t.GetPrimaryKeyColumns()
	nonPkColumns := t.GetNonPrimaryKeyColumns()

	if len(pkColumns) > 0 {
		for _, column := range pkColumns {
			str := ""
			if column.NotNull {
				str = fmt.Sprintf("  * %s : %s", column.Name, column.Type)
			} else {
				str = fmt.Sprintf("  %s : %s", column.Name, column.Type)
			}
			lines = append(lines, str)
		}
	}

	if len(pkColumns) > 0 && len(nonPkColumns) > 0 {
		lines = append(lines, "  --")
	}

	if len(nonPkColumns) > 0 {
		for _, column := range nonPkColumns {
			str := ""
			if column.NotNull {
				str = fmt.Sprintf("  * %s : %s", column.Name, column.Type)
			} else {
				str = fmt.Sprintf("  %s : %s", column.Name, column.Type)
			}
			lines = append(lines, str)
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// GetPrimaryKeyColumns returns Primary key columns
func (t *Table) GetPrimaryKeyColumns() []Column {
	var columns []Column
	for _, column := range t.Columns {
		if column.PrimaryKey {
			columns = append(columns, column)
		}
	}
	return columns
}

// GetNonPrimaryKeyColumns returns Non-Primary key columns
func (t *Table) GetNonPrimaryKeyColumns() []Column {
	var columns []Column
	for _, column := range t.Columns {
		if !column.PrimaryKey {
			columns = append(columns, column)
		}
	}
	return columns
}
