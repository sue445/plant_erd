package adapter

// ForeignKey represents foreign key info
type ForeignKey struct {
	Sequence   int
	FromColumn string
	ToTable    string
	ToColumn   string
}
