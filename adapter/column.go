package adapter

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
