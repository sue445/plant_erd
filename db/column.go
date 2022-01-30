package db

import "fmt"

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
	return fmt.Sprintf("%s %s", c.Type, c.Name)
}
