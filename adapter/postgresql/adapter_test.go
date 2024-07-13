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
	adapter, closeDatabase, err := NewAdapter(config)

	if err != nil {
		panic(err)
	}

	defer closeDatabase() //nolint:errcheck

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

		a.db.MustExec(`CREATE SCHEMA people;`)
		defer func() {
			a.db.MustExec("DROP SCHEMA people;")
		}()

		a.db.MustExec(`
			CREATE TABLE people.author (
				id int4 NOT NULL,
				"name" varchar NOT NULL,
				CONSTRAINT author_pk PRIMARY KEY (id)
			);
		`)
		defer func() {
			a.db.MustExec("DROP TABLE people.author;")
		}()

		a.db.MustExec(`
			CREATE TABLE people.book_author (
				isbn varchar NOT NULL,
				author_id int4 NOT NULL,
				CONSTRAINT book_author_pk PRIMARY KEY (isbn, author_id),
				CONSTRAINT book_author_fk_1 FOREIGN KEY (author_id) REFERENCES people.author(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
		defer func() {
			a.db.MustExec("DROP TABLE people.book_author;")
		}()

		a.db.MustExec(`CREATE SCHEMA "library";`)
		defer func() {
			a.db.MustExec(`DROP SCHEMA "library";`)
		}()

		a.db.MustExec(`
			CREATE TABLE "library".book (
				isbn varchar NOT NULL,
				"name" varchar NOT NULL,
				CONSTRAINT book_pk PRIMARY KEY (isbn)
			);
		`)
		defer func() {
			a.db.MustExec(`DROP TABLE "library".book;`)
		}()

		gotTables, err := a.GetAllTableNames()

		wantTables := []string{
			"library.book",
			"people.author",
			"people.book_author",
			"public.articles",
			"public.users",
		}

		if assert.NoError(t, err) {
			assert.Equal(t, wantTables, gotTables)
		}
	})
}

func TestAdapter_GetTable_in_public_schema(t *testing.T) {
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
				name: "public.users",
				args: args{
					tableName: "public.users",
				},
				want: &db.Table{
					Name: "public.users",
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
				name: "public.articles",
				args: args{
					tableName: "public.articles",
				},
				want: &db.Table{
					Name: "public.articles",
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
							ToTable:    "public.users",
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
				name: "public.followers",
				args: args{
					tableName: "public.followers",
				},
				want: &db.Table{
					Name: "public.followers",
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
							ToTable:    "public.users",
							ToColumn:   "id",
						},
						{
							FromColumn: "user_id",
							ToTable:    "public.users",
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

func TestAdapter_GetTable_in_non_public_schema(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`CREATE SCHEMA people;`)
		defer func() {
			a.db.MustExec("DROP SCHEMA people;")
		}()

		a.db.MustExec(`
			CREATE TABLE people.author (
				id int4 NOT NULL,
				"name" varchar NOT NULL,
				CONSTRAINT author_pk PRIMARY KEY (id)
			);
		`)
		defer func() {
			a.db.MustExec("DROP TABLE people.author;")
		}()

		a.db.MustExec(`
			CREATE TABLE people.book_author (
				isbn varchar NOT NULL,
				author_id int4 NOT NULL,
				CONSTRAINT book_author_pk PRIMARY KEY (isbn, author_id),
				CONSTRAINT book_author_fk_1 FOREIGN KEY (author_id) REFERENCES people.author(id) ON DELETE CASCADE ON UPDATE CASCADE
			);
		`)
		defer func() {
			a.db.MustExec("DROP TABLE people.book_author;")
		}()

		a.db.MustExec(`CREATE SCHEMA "library";`)
		defer func() {
			a.db.MustExec(`DROP SCHEMA "library";`)
		}()

		a.db.MustExec(`
			CREATE TABLE "library".book (
				isbn varchar NOT NULL,
				"name" varchar NOT NULL,
				CONSTRAINT book_pk PRIMARY KEY (isbn)
			);
		`)
		defer func() {
			a.db.MustExec(`DROP TABLE "library".book;`)
		}()

		type args struct {
			tableName string
		}
		tests := []struct {
			name string
			args args
			want *db.Table
		}{
			{
				name: "people.author",
				args: args{
					tableName: "people.author",
				},
				want: &db.Table{
					Name: "people.author",
					Columns: []*db.Column{
						{
							Name:       "id",
							Type:       "integer",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "name",
							Type:    "character varying",
							NotNull: true,
						},
					},
				},
			},
			{
				name: "people.book_author",
				args: args{
					tableName: "people.book_author",
				},
				want: &db.Table{
					Name: "people.book_author",
					Columns: []*db.Column{
						{
							Name:       "isbn",
							Type:       "character varying",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:       "author_id",
							Type:       "integer",
							NotNull:    true,
							PrimaryKey: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							FromColumn: "author_id",
							ToTable:    "people.author",
							ToColumn:   "id",
						},
					},
				},
			},
			{
				name: "library.book",
				args: args{
					tableName: "library.book",
				},
				want: &db.Table{
					Name: "library.book",
					Columns: []*db.Column{
						{
							Name:       "isbn",
							Type:       "character varying",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "name",
							Type:    "character varying",
							NotNull: true,
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
