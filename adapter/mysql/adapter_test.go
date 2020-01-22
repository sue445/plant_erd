package mysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
	"os"
	"regexp"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("MYSQL_HOST") == "" || os.Getenv("MYSQL_PORT") == "" || os.Getenv("MYSQL_USER") == "" || os.Getenv("MYSQL_DATABASE") == "" {
		println("adapter/mysql test is skipped because MYSQL_HOST, MYSQL_PORT and MYSQL_USER not all set")
		return
	}

	ret := m.Run()
	os.Exit(ret)
}

func withDatabase(callback func(*Adapter)) {
	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT")
	config.User = os.Getenv("MYSQL_USER")
	config.Passwd = os.Getenv("MYSQL_PASSWORD")
	config.DBName = os.Getenv("MYSQL_DATABASE")
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
				id   int not null primary key, 
				name varchar(191)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE users;")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      int not null primary key, 
				user_id int not null, 
				FOREIGN KEY fk_users (user_id) REFERENCES users(id)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE articles;")
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
				id   int not null primary key, 
				name varchar(191)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE users;")
		}()

		a.db.MustExec(`
			CREATE TABLE articles (
				id      int not null primary key, 
				user_id int not null, 
				FOREIGN KEY fk_users (user_id) REFERENCES users(id)
		);`)
		defer func() {
			a.db.MustExec("DROP TABLE articles;")
		}()
		a.db.MustExec("CREATE INDEX index_user_id_on_articles ON articles(user_id)")

		a.db.MustExec(`
			CREATE TABLE followers (
				id             int not null primary key,
				user_id        int not null,
				target_user_id int not null,
				FOREIGN KEY fk_users (user_id)         REFERENCES users(id),
				FOREIGN KEY fk_users2 (target_user_id) REFERENCES users(id)
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
							Type:       "int",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name: "name",
							Type: "varchar(191)",
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
							Type:       "int",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "user_id",
							Type:    "int",
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
							Type:       "int",
							NotNull:    true,
							PrimaryKey: true,
						},
						{
							Name:    "user_id",
							Type:    "int",
							NotNull: true,
						},
						{
							Name:    "target_user_id",
							Type:    "int",
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
							Name:    "index_user_id_and_target_user_id_on_followers",
							Columns: []string{"user_id", "target_user_id"},
							Unique:  true,
						},
						{
							Name:    "index_target_user_id_and_user_id_on_followers",
							Columns: []string{"target_user_id", "user_id"},
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
					// assert.Equal(t, tt.want, got)

					assert.Equal(t, tt.want.Name, got.Name)
					assert.Equal(t, tt.want.ForeignKeys, got.ForeignKeys)
					assert.Equal(t, tt.want.Indexes, got.Indexes)

					if assert.Equal(t, len(tt.want.Columns), len(got.Columns)) {
						for i := 0; i < len(got.Columns); i++ {
							wantColumn := tt.want.Columns[i]
							gotColumn := got.Columns[i]

							assert.Equal(t, wantColumn.Name, gotColumn.Name)
							assert.Equal(t, wantColumn.NotNull, gotColumn.NotNull)
							assert.Equal(t, wantColumn.PrimaryKey, gotColumn.PrimaryKey)

							// FIXME: Type is `int(11)` when MySQL 5.6 and 5.7, but Type is `int` when MySQL 8
							if wantColumn.Type == "int" {
								assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("^%s(\\(\\d+\\))?$", wantColumn.Type)), gotColumn.Type)
							} else {
								assert.Equal(t, wantColumn.Type, gotColumn.Type)
							}
						}
					}
				}
			})
		}
	})
}
