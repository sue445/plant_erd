package postgresql

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
)

func TestMain(m *testing.M) {
	if os.Getenv("POSTGRES_HOST") == "" || os.Getenv("POSTGRES_PORT") == "" || os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_DATABASE") == "" {
		println("adapter/postgresql test is skipped because POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER and POSTGRES_DATABASE not all set")
		return
	}

	ret := m.Run()
	os.Exit(ret)
}

func withDatabase(callback func(*Adapter)) {
	config := NewConfig()
	config.Host = os.Getenv("POSTGRES_HOST")
	config.Port, _ = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	config.User = os.Getenv("POSTGRES_USER")
	config.Password = os.Getenv("POSTGRES_PASSWORD")
	config.DBName = os.Getenv("POSTGRES_DATABASE")
	adapter, close, err := NewAdapter(config)

	if err != nil {
		panic(err)
	}

	defer close()

	adapter.db.MustExec("DROP TABLE IF EXISTS followers;")
	adapter.db.MustExec("DROP TABLE IF EXISTS articles;")
	adapter.db.MustExec("DROP TABLE IF EXISTS users;")
	callback(adapter)
}

func TestAdapter_GetAllTableNames(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name varchar(191)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE users;")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY (user_id) REFERENCES users(id)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE articles;")
		}()

		a.db.MustExec("CREATE VIEW user_names AS SELECT name FROM users;")
		defer func() {
			a.db.MustExec("DROP VIEW user_names;")
		}()

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
		defer func() {
			a.db.MustExec("DROP TABLE users;")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY(user_id) REFERENCES users(id)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE articles;")
		}()
		a.db.MustExec("CREATE INDEX index_user_id_on_articles ON articles(user_id)")

		a.db.MustExec(`
			CREATE TABLE followers (
				id             integer not null primary key,
				user_id        integer not null,
				target_user_id integer not null,
				FOREIGN KEY(user_id)        REFERENCES users(id),
				FOREIGN KEY(target_user_id) REFERENCES users(id)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE followers;")
		}()
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
