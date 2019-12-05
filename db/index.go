package db

// Index represents index definition
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}
