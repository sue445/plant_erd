package db

import (
	"fmt"
	"strings"
)

// Column represents column info
type Column struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
}

// ToErd returns ERD formatted column
func (c *Column) ToErd() string {
	str := ""

	if c.NotNull {
		str += "* "
	}

	str += fmt.Sprintf("%s : %s", c.Name, c.Type)

	return str
}

// ToMermaid returns Mermaid formatted column
func (c *Column) ToMermaid() string {
	mermaidType := c.Type

	// mermaid cannot display Type Column "()" and "unsigned"
	mermaidType = strings.Replace(mermaidType, "(", "[", -1)
	mermaidType = strings.Replace(mermaidType, ")", "]", -1)

	mermaidType = strings.Replace(mermaidType, " ", "_", -1)

	return fmt.Sprintf("%s %s", mermaidType, c.Name)
}
