package oracle

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
)

func TestMain(m *testing.M) {
	if os.Getenv("ORACLE_HOST") == "" || os.Getenv("ORACLE_PORT") == "" || os.Getenv("ORACLE_USER") == "" || os.Getenv("ORACLE_SERVICE") == "" {
		println("adapter/oracle test is skipped because ORACLE_HOST, ORACLE_PORT, ORACLE_USER and ORACLE_SERVICE not all set")
		return
	}

	ret := m.Run()
	os.Exit(ret)
}

func withDatabase(callback func(*Adapter)) {
	config := NewConfig()
	config.Host = os.Getenv("ORACLE_HOST")
	config.Port, _ = strconv.Atoi(os.Getenv("ORACLE_PORT"))
	config.Username = os.Getenv("ORACLE_USER")
	config.Password = os.Getenv("ORACLE_PASSWORD")
	config.ServiceName = os.Getenv("ORACLE_SERVICE")
	adapter, closeDatabase, err := NewAdapter(config)

	if err != nil {
		panic(err)
	}

	defer closeDatabase() //nolint:errcheck

	// adapter.db.MustExec("DROP TABLE followers")
	// adapter.db.MustExec("DROP TABLE articles")
	// adapter.db.MustExec("DROP TABLE users")
	callback(adapter)
}

func TestAdapter_GetAllTableNames(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name varchar(191)
		)`)
		defer func() {
			a.db.MustExec("DROP TABLE users")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY (user_id) REFERENCES users(id)
		)`)
		defer func() {
			a.db.MustExec("DROP TABLE articles")
		}()

		a.db.MustExec("CREATE VIEW user_names AS SELECT name FROM users")
		defer func() {
			a.db.MustExec("DROP VIEW user_names")
		}()

		tables, err := a.GetAllTableNames()

		if assert.NoError(t, err) {
			assert.Contains(t, tables, "articles")
			assert.Contains(t, tables, "users")
		}
	})
}

func TestAdapter_GetTable(t *testing.T) {
	withDatabase(func(a *Adapter) {
		a.db.MustExec(`
			CREATE TABLE users (
				id   integer not null primary key,
				name varchar2(191)
		)`)
		defer func() {
			a.db.MustExec("DROP TABLE users")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      integer not null primary key,
				user_id integer not null,
				FOREIGN KEY(user_id) REFERENCES users(id)
		)`)
		defer func() {
			a.db.MustExec("DROP TABLE articles")
		}()
		a.db.MustExec("CREATE INDEX user_id ON articles(user_id)")

		a.db.MustExec(`
			CREATE TABLE followers (
				id             integer not null primary key,
				user_id        integer not null,
				target_user_id integer not null,
				FOREIGN KEY(user_id)        REFERENCES users(id),
				FOREIGN KEY(target_user_id) REFERENCES users(id)
		)`)
		defer func() {
			a.db.MustExec("DROP TABLE followers")
		}()
		a.db.MustExec("CREATE UNIQUE INDEX user_id_target_user_id ON followers(user_id, target_user_id)")
		a.db.MustExec("CREATE UNIQUE INDEX target_user_id_user_id ON followers(target_user_id, user_id)")

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
							Name:       "ID",
							Type:       "NUMBER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name: "NAME",
							Type: "VARCHAR2(191)",
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
							Name:       "ID",
							Type:       "NUMBER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "USER_ID",
							Type:    "NUMBER",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							FromColumn: "USER_ID",
							ToTable:    "USERS",
							ToColumn:   "ID",
						},
					},
					Indexes: []*db.Index{
						{
							Name:    "USER_ID",
							Columns: []string{"USER_ID"},
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
							Name:       "ID",
							Type:       "NUMBER",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "USER_ID",
							Type:    "NUMBER",
							NotNull: true,
						},
						{
							Name:    "TARGET_USER_ID",
							Type:    "NUMBER",
							NotNull: true,
						},
					},
					ForeignKeys: []*db.ForeignKey{
						{
							FromColumn: "TARGET_USER_ID",
							ToTable:    "USERS",
							ToColumn:   "ID",
						},
						{
							FromColumn: "USER_ID",
							ToTable:    "USERS",
							ToColumn:   "ID",
						},
					},
					Indexes: []*db.Index{
						{
							Name:    "USER_ID_TARGET_USER_ID",
							Columns: []string{"USER_ID", "TARGET_USER_ID"},
							Unique:  true,
						},
						{
							Name:    "TARGET_USER_ID_USER_ID",
							Columns: []string{"TARGET_USER_ID", "USER_ID"},
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
