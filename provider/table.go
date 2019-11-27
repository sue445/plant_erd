package provider

// Table represents table info
type Table struct {
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
}

// Column represents column info
type Column struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
}

// ForeignKey represents foreign key info
type ForeignKey struct {
	Sequence   int
	FromColumn string
	ToTable    string
	ToColumn   string
}
