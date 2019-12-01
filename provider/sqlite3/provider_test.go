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
		CREATE TABLE articles (id integer not null primary key, user_id integer not null, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		p.db.MustExec(sql)

		tables, err := p.GetAllTableNames()
		assert.NoError(t, err)

		if err == nil {
			assert.Equal(t, []string{"articles", "users"}, tables)
		}
	})
}

func TestProvider_GetTable(t *testing.T) {
	withDatabase(func(p *Provider) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer not null, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		p.db.MustExec(sql)

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
					Columns: []*provider.Column{
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
					Columns: []*provider.Column{
						{
							Name:       "id",
							Type:       "integer",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "user_id",
							Type:    "integer",
							NotNull: true,
						},
					},
					ForeignKeys: []*provider.ForeignKey{
						{
							Sequence:   0,
							FromColumn: "user_id",
							ToTable:    "users",
							ToColumn:   "id",
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := p.GetTable(tt.args.tableName)

				assert.NoError(t, err)

				if err == nil {
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}
