package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func withDatabase(callback func(*Provider)) {
	provider, close, err := NewProvider("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	defer close()

	callback(provider)
}

func TestProvider_GetAllTableNames(t *testing.T) {
	withDatabase(func(provider *Provider) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		_, err := provider.db.Exec(sql)
		if err != nil {
			panic(err)
		}

		tables, err := provider.GetAllTableNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"articles", "users"}, tables)
	})
}
