package postgresql

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("POSTGRES_HOST") == "" || os.Getenv("POSTGRES_PORT") == "" || os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_DATABASE") == "" {
		println("adapter/postgresql test is skipped because POSTGRES_HOST, POSTGRES_PORT and POSTGRES_USER not all set")
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
				FOREIGN KEY (user_id) REFERENCES users(id)
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
