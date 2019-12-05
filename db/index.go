package db

import (
	"fmt"
	"strings"
)

// Index represents index definition
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

// ToErd returns ERD formatted index
func (i *Index) ToErd() string {
	str := ""

	if i.Unique {
		str += "- "
	}

	str += fmt.Sprintf("%s (%s)", i.Name, strings.Join(i.Columns, ", "))

	return str
}
