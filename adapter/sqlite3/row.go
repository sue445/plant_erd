package sqlite3

type sqliteMaster struct {
	Name string `db:"name"`
}

type tableInfo struct {
	CID        int     `db:"cid"`
	Name       string  `db:"name"`
	Type       string  `db:"type"`
	NotNull    bool    `db:"notnull"`
	DfltValue  *string `db:"dflt_value"`
	PrimaryKey bool    `db:"pk"`
}

type foreignKeyList struct {
	ID       int    `db:"id"`
	Seq      int    `db:"seq"`
	Table    string `db:"table"`
	From     string `db:"from"`
	To       string `db:"to"`
	OnUpdate string `db:"on_update"`
	OnDelete string `db:"on_delete"`
	Match    string `db:"match"`
}
