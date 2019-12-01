package provider

// Column represents column info
type Column struct {
	Name       string
	Type       string
	NotNull    bool
	PrimaryKey bool
}
