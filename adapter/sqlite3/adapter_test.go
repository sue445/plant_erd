package sqlite3

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
	"testing"
)

func withDatabase(callback func(*Adapter)) {
	adapter, close, err := NewAdapter("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	defer close()

	adapter.db.MustExec("PRAGMA foreign_keys = ON;")

	callback(adapter)
}

func TestAdapter_GetAllTableNames(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name text
		);`)

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key, 
				user_id integer not null, 
				FOREIGN KEY(user_id) REFERENCES users(id)
		);`)

		tables, err := a.GetAllTableNames()

		if assert.NoError(t, err) {
			assert.Equal(t, []string{"articles", "users"}, tables)
		}
	})
}

func TestAdapter_GetTable(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name text
		);`)

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key, 
				user_id integer not null, 
				FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
		a.db.MustExec("CREATE INDEX index_user_id_on_articles ON articles(user_id)")

		a.db.MustExec(`
			CREATE TABLE followers (
				id             integer not null primary key,
				user_id        integer not null,
				target_user_id integer not null,
				FOREIGN KEY(user_id)        REFERENCES users(id),
				FOREIGN KEY(target_user_id) REFERENCES users(id)
		);`)
		a.db.MustExec("CREATE UNIQUE INDEX index_user_id_and_target_user_id_on_followers ON followers(user_id, target_user_id)")
		a.db.MustExec("CREATE UNIQUE INDEX index_target_user_id_and_user_id_on_followers ON followers(target_user_id, user_id)")

		type args struct {
			tableName string
		}
		tests := []struct {
			name string
			args args
			want *db.Table
		}{
			{
				name: "users",
				args: args{
					tableName: "users",
				},
				want: &db.Table{
					Name: "users",
					Columns: []*db.Column{
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
				want: &db.Table{
					Name: "articles",
					Columns: []*db.Column{
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
					ForeignKeys: []*db.ForeignKey{
						{
							Sequence:   0,
							FromColumn: "user_id",
							ToTable:    "users",
							ToColumn:   "id",
						},
					},
					Indexes: []*db.Index{
						{
							Name:    "index_user_id_on_articles",
							Columns: []string{"user_id"},
							Unique:  false,
						},
					},
				},
			},
			{
				name: "followers",
				args: args{
					tableName: "followers",
				},
				want: &db.Table{
					Name: "followers",
					Columns: []*db.Column{
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
						{
							Name:    "target_user_id",
							Type:    "integer",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							Sequence:   0,
							FromColumn: "target_user_id",
							ToTable:    "users",
							ToColumn:   "id",
						},
						{
							Sequence:   0,
							FromColumn: "user_id",
							ToTable:    "users",
							ToColumn:   "id",
						},
					},
					Indexes: []*db.Index{
						{
							Name:    "index_target_user_id_and_user_id_on_followers",
							Columns: []string{"target_user_id", "user_id"},
							Unique:  true,
						},
						{
							Name:    "index_user_id_and_target_user_id_on_followers",
							Columns: []string{"user_id", "target_user_id"},
							Unique:  true,
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := a.GetTable(tt.args.tableName)

				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}
