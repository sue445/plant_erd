package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/provider"
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
	withDatabase(func(p *Provider) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		_, err := p.db.Exec(sql)
		if err != nil {
			panic(err)
		}

		tables, err := p.GetAllTableNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"articles", "users"}, tables)
	})
}

func TestProvider_GetTable(t *testing.T) {
	withDatabase(func(p *Provider) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		_, err := p.db.Exec(sql)
		if err != nil {
			panic(err)
		}

		type args struct {
			tableName string
		}
		tests := []struct {
			name string
			args args
			want *provider.Table
		}{
			{
				name: "users",
				args: args{
					tableName: "users",
				},
				want: &provider.Table{
					Name: "users",
					Columns: []provider.Column{
						{
							Name:       "id",
							Type:       "integer",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name: "name",
							Type: "text",
						},
					},
				},
			},
			{
				name: "articles",
				args: args{
					tableName: "articles",
				},
				want: &provider.Table{
					Name: "articles",
					Columns: []provider.Column{
						{
							Name:       "id",
							Type:       "integer",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name: "user_id",
							Type: "integer",
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := p.GetTable(tt.args.tableName)

				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
