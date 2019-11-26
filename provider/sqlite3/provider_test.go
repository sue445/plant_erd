package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestDb() (*Provider, Close) {
	provider, close, err := NewProvider("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	return provider, close
}

func TestProvider_GetTableNames(t *testing.T) {
	provider, close := setupTestDb()
	defer close()

	sql := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE users (id integer not null primary key, name text);
	CREATE TABLE articles (id integer not null primary key, user_id integer, FOREIGN KEY(user_id) REFERENCES users(id));
	`

	_, err := provider.db.Exec(sql)
	if err != nil {
		panic(err)
	}

	tables, err := provider.GetTableNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"articles", "users"}, tables)
}
