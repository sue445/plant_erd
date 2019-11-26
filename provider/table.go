package provider

// Table represents table info
type Table struct {
	Name    string
	Columns []Column
}

// Column represents column info
type Column struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
}
