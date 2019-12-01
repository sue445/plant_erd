package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/adapter"
	"testing"
)

func withDatabase(callback func(*Adapter)) {
	adapter, close, err := NewAdapter("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	defer close()

	callback(adapter)
}

func TestAdapter_GetAllTableNames(t *testing.T) {
	withDatabase(func(a *Adapter) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer not null, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		a.db.MustExec(sql)

		tables, err := a.GetAllTableNames()
		assert.NoError(t, err)

		if err == nil {
			assert.Equal(t, []string{"articles", "users"}, tables)
		}
	})
}

func TestAdapter_GetTable(t *testing.T) {
	withDatabase(func(a *Adapter) {
		sql := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE users (id integer not null primary key, name text);
		CREATE TABLE articles (id integer not null primary key, user_id integer not null, FOREIGN KEY(user_id) REFERENCES users(id));
		`

		a.db.MustExec(sql)

		type args struct {
			tableName string
		}
		tests := []struct {
			name string
			args args
			want *adapter.Table
		}{
			{
				name: "users",
				args: args{
					tableName: "users",
				},
				want: &adapter.Table{
					Name: "users",
					Columns: []*adapter.Column{
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
				want: &adapter.Table{
					Name: "articles",
					Columns: []*adapter.Column{
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
					ForeignKeys: []*adapter.ForeignKey{
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
				got, err := a.GetTable(tt.args.tableName)

				assert.NoError(t, err)

				if err == nil {
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}
