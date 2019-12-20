package oracle

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("ORACLE_HOST") == "" || os.Getenv("ORACLE_PORT") == "" || os.Getenv("ORACLE_USER") == "" || os.Getenv("ORACLE_SERVICE") == "" {
		println("adapter/oracle test is skipped because ORACLE_HOST, ORACLE_PORT and ORACLE_USER and ORACLE_SERVICE not all set")
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
	adapter, close, err := NewAdapter(config)

	if err != nil {
		panic(err)
	}

	defer close()

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

		tables, err := a.GetAllTableNames()

		if assert.NoError(t, err) {
			assert.Contains(t, tables, "articles")
			assert.Contains(t, tables, "users")
		}
	})
}
