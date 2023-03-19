package sqlite3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
)

func withDatabase(callback func(*Adapter)) {
	adapter, close, err := NewAdapter("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	defer close()

	adapter.DB.MustExec("PRAGMA foreign_keys = ON;")

	callback(adapter)
}

func TestAdapter_GetAllTableNames(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.DB.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name text
		);`)

		a.DB.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY(user_id) REFERENCES users(id)
		);`)

		a.DB.MustExec("CREATE VIEW user_names AS SELECT name FROM users;")

		tables, err := a.GetAllTableNames()

		if assert.NoError(t, err) {
			assert.Equal(t, []string{"articles", "users"}, tables)
		}
	})
}

func TestAdapter_GetTable(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.DB.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name text
		);`)

		a.DB.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
		a.DB.MustExec("CREATE INDEX index_user_id_on_articles ON articles(user_id)")

		a.DB.MustExec(`
			CREATE TABLE followers (
				id             integer not null primary key,
				user_id        integer not null,
				target_user_id integer not null,
				FOREIGN KEY(user_id)        REFERENCES users(id),
				FOREIGN KEY(target_user_id) REFERENCES users(id)
		);`)
		a.DB.MustExec("CREATE UNIQUE INDEX index_user_id_and_target_user_id_on_followers ON followers(user_id, target_user_id)")
		a.DB.MustExec("CREATE UNIQUE INDEX index_target_user_id_and_user_id_on_followers ON followers(target_user_id, user_id)")

		a.DB.MustExec(`
			CREATE TABLE album_genres (
				album_id varchar default null not null
					references album
						on delete cascade,
				genre_id varchar default null not null
					references genre
						on delete cascade,
				constraint album_genre_ux
					unique (album_id, genre_id)
		);`)

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
							Type:       "INTEGER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name: "name",
							Type: "TEXT",
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
							Type:       "INTEGER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "user_id",
							Type:    "INTEGER",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
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
							Type:       "INTEGER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "user_id",
							Type:    "INTEGER",
							NotNull: true,
						},
						{
							Name:    "target_user_id",
							Type:    "INTEGER",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							FromColumn: "target_user_id",
							ToTable:    "users",
							ToColumn:   "id",
						},
						{
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
			{
				name: "album_genres",
				args: args{
					tableName: "album_genres",
				},
				want: &db.Table{
					Name: "album_genres",
					Columns: []*db.Column{
						{
							Name:    "album_id",
							Type:    "varchar",
							NotNull: true,
						},
						{
							Name:    "genre_id",
							Type:    "varchar",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							FromColumn: "genre_id",
							ToTable:    "genre",
							ToColumn:   "id",
						},
						{
							FromColumn: "album_id",
							ToTable:    "album",
							ToColumn:   "id",
						},
					},
					Indexes: []*db.Index{
						{
							Name:    "sqlite_autoindex_album_genres_1",
							Columns: []string{"album_id", "genre_id"},
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
