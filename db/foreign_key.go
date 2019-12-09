package db

// ForeignKey represents foreign key info
type ForeignKey struct {
	FromColumn string
	ToTable    string
	ToColumn   string
}
