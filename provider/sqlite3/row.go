package sqlite3

type sqliteMasterRow struct {
	Name string `db:"name"`
}

type pragmaTableInfoRow struct {
	CID        int     `db:"cid"`
	Name       string  `db:"name"`
	Type       string  `db:"type"`
	NotNull    bool    `db:"notnull"`
	DfltValue  *string `db:"dflt_value"`
	PrimaryKey bool    `db:"pk"`
}

type pragmaForeignKeyListRow struct {
	ID       int    `db:"id"`
	Seq      int    `db:"seq"`
	Table    string `db:"table"`
	From     string `db:"from"`
	To       string `db:"to"`
	OnUpdate string `db:"on_update"`
	OnDelete string `db:"on_delete"`
	Match    string `db:"match"`
}
