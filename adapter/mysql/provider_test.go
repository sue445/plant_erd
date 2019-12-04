package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"os"
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
	config := &mysql.Config{
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		DBName: os.Getenv("MYSQL_DATABASE"),
	}
	adapter, close, err := NewAdapter(config)

	if err != nil {
		panic(err)
	}

	defer close()

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
