package lib

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/adapter/sqlite3"
	"github.com/sue445/plant_erd/db"
)

func withDatabase(callback func(*sqlite3.Adapter)) {
	adapter, close, err := sqlite3.NewAdapter("file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}

	defer close()

	adapter.DB.MustExec("PRAGMA foreign_keys = ON;")

	callback(adapter)
}

func TestErdGenerator_generate(t *testing.T) {
	tables := []*db.Table{
		{
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
		},
		{
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
	}
	schema := db.NewSchema(tables)

	type fields struct {
		Filepath string
		Table    string
		Distance int
	}
	type args struct {
		schema *db.Schema
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "no table",
			fields: fields{
				Table:    "",
				Distance: 0,
			},
			args: args{
				schema: schema,
			},
		},
		{
			name: "with table and distance",
			fields: fields{
				Table:    "users",
				Distance: 1,
			},
			args: args{
				schema: schema,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &ErdGenerator{
				Filepath: tt.fields.Filepath,
				Table:    tt.fields.Table,
				Distance: tt.fields.Distance,
			}
			got := g.generate(tt.args.schema)
			assert.Greater(t, len(got), 0)
		})
	}
}

func TestErdGenerator_outputErd_ToFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir)

	filePath := filepath.Join(dir, "erd.txt")
	g := &ErdGenerator{
		Filepath: filePath,
	}

	g.outputErd("aaa")

	data, err := ioutil.ReadFile(filePath)

	if assert.NoError(t, err) {
		str := string(data)
		assert.Equal(t, "aaa", str)
	}
}

// c.f. https://qiita.com/kami_zh/items/ff636f15da87dabebe6c
func captureStdout(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestErdGenerator_outputErd_ToStdout(t *testing.T) {
	g := &ErdGenerator{
		Filepath: "",
	}

	str := captureStdout(func() {
		err := g.outputErd("aaa")
		assert.NoError(t, err)
	})

	assert.Equal(t, "aaa", str)
}

func TestErdGenerator_generate_withSkipTable(t *testing.T) {
	tables := []*db.Table{
		{
			Name: "articles",
		},
		{
			Name: "users",
		},
		{
			Name: "QRTZ_TRIGGERS",
		},
		{
			Name: "QRTZ_ALARMS",
		},
		{
			Name: "QRTZ_SCHEDULER",
		},
	}
	schema := db.NewSchema(tables)

	type fields struct {
		SkipTable string
	}
	type args struct {
		schema *db.Schema
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "with skip tables begin with QRTZ*",
			fields: fields{
				SkipTable: "(QRTZ*)\\w+",
			},
			args: args{
				schema: schema,
			},
		},
		{
			name: "with skip all tables",
			fields: fields{
				SkipTable: "()\\w+",
			},
			args: args{
				schema: schema,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &ErdGenerator{
				SkipTable: tt.fields.SkipTable,
			}
			got := g.generate(tt.args.schema)
			if tt.name == "with skip tables begin with QRTZ*" {
				assert.Contains(t, got, "articles", "users")
				assert.NotContains(t, got, "QRTZ_TRIGGERS", "QRTZ_ALARMS", "QRTZ_SCHEDULER")
			}
			if tt.name == "with skip all tables" {
				assert.Equal(t, len(got), 0)
				assert.NotContains(t, got, "articles", "users", "QRTZ_TRIGGERS", "QRTZ_ALARMS", "QRTZ_SCHEDULER")
			}
		})
	}
}

func TestErdGenerator_checkParamTable(t *testing.T) {
	type fields struct {
		Filepath string
		Table    string
		Distance int
	}
	type args struct {
		schema *db.Schema
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "--table is not passed",
			fields: fields{
				Table: "",
			},
			args: args{
				schema: &db.Schema{
					Tables: []*db.Table{
						{
							Name: "users",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "--table is passed and tables is exists",
			fields: fields{
				Table: "users",
			},
			args: args{
				schema: &db.Schema{
					Tables: []*db.Table{
						{
							Name: "users",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "--table is passed and tables is not exists",
			fields: fields{
				Table: "users",
			},
			args: args{
				schema: &db.Schema{
					Tables: []*db.Table{
						{
							Name: "articles",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &ErdGenerator{
				Filepath: tt.fields.Filepath,
				Table:    tt.fields.Table,
				Distance: tt.fields.Distance,
			}

			err := g.checkParamTable(tt.args.schema)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func ExampleErdGenerator_Run_two_tables() {
	withDatabase(func(a *sqlite3.Adapter) {
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

		schema, err := LoadSchema(a)
		if err != nil {
			panic(err)
		}

		generator := ErdGenerator{}
		generator.Run(schema)

		// Output:
		// entity articles {
		//   * id : integer
		//   --
		//   * user_id : integer
		//   --
		//   index_user_id_on_articles (user_id)
		// }
		//
		// entity users {
		//   * id : integer
		//   --
		//   name : text
		// }
		//
		// articles }-- users
	})
}

func createManyExampleTables(a *sqlite3.Adapter) {
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
		CREATE TABLE comments (
			id         integer not null primary key, 
			article_id integer not null, 
			FOREIGN KEY(article_id) REFERENCES articles(id)
	);`)
	a.DB.MustExec("CREATE INDEX index_article_id_on_articles ON comments(article_id)")

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
		CREATE TABLE followings (
			id             integer not null primary key,
			user_id        integer not null,
			target_user_id integer not null,
			FOREIGN KEY(user_id)        REFERENCES users(id),
			FOREIGN KEY(target_user_id) REFERENCES users(id)
	);`)
	a.DB.MustExec("CREATE UNIQUE INDEX index_user_id_and_target_user_id_on_followings ON followings(user_id, target_user_id)")
	a.DB.MustExec("CREATE UNIQUE INDEX index_target_user_id_and_user_id_on_followings ON followings(target_user_id, user_id)")

	a.DB.MustExec(`
		CREATE TABLE likes (
			article_id integer not null, 
			user_id    integer not null, 
			FOREIGN KEY(article_id) REFERENCES articles(id)
			FOREIGN KEY(user_id)    REFERENCES users(id)
	);`)
	a.DB.MustExec("CREATE UNIQUE INDEX index_article_id_and_user_id_on_likes ON likes(article_id, user_id)")
	a.DB.MustExec("CREATE INDEX index_user_id_on_likes ON likes(user_id)")

	a.DB.MustExec(`
		CREATE TABLE revisions (
			id         integer not null primary key, 
			article_id integer not null, 
			FOREIGN KEY(article_id) REFERENCES articles(id)
	);`)
	a.DB.MustExec("CREATE INDEX index_article_id_on_revisions ON revisions(article_id)")
}

func ExampleErdGenerator_Run_many_tables() {
	withDatabase(func(a *sqlite3.Adapter) {
		createManyExampleTables(a)

		schema, err := LoadSchema(a)
		if err != nil {
			panic(err)
		}

		generator := ErdGenerator{}
		generator.Run(schema)

		// Output:
		// entity articles {
		//   * id : integer
		//   --
		//   * user_id : integer
		//   --
		//   index_user_id_on_articles (user_id)
		// }
		//
		// entity comments {
		//   * id : integer
		//   --
		//   * article_id : integer
		//   --
		//   index_article_id_on_articles (article_id)
		// }
		//
		// entity followers {
		//   * id : integer
		//   --
		//   * user_id : integer
		//   * target_user_id : integer
		//   --
		//   - index_target_user_id_and_user_id_on_followers (target_user_id, user_id)
		//   - index_user_id_and_target_user_id_on_followers (user_id, target_user_id)
		// }
		//
		// entity followings {
		//   * id : integer
		//   --
		//   * user_id : integer
		//   * target_user_id : integer
		//   --
		//   - index_target_user_id_and_user_id_on_followings (target_user_id, user_id)
		//   - index_user_id_and_target_user_id_on_followings (user_id, target_user_id)
		// }
		//
		// entity likes {
		//   * article_id : integer
		//   * user_id : integer
		//   --
		//   index_user_id_on_likes (user_id)
		//   - index_article_id_and_user_id_on_likes (article_id, user_id)
		// }
		//
		// entity revisions {
		//   * id : integer
		//   --
		//   * article_id : integer
		//   --
		//   index_article_id_on_revisions (article_id)
		// }
		//
		// entity users {
		//   * id : integer
		//   --
		//   name : text
		// }
		//
		// articles }-- users
		//
		// comments }-- articles
		//
		// followers }-- users
		//
		// followers }-- users
		//
		// followings }-- users
		//
		// followings }-- users
		//
		// likes }-- users
		//
		// likes }-- articles
		//
		// revisions }-- articles
	})
}

func ExampleErdGenerator_Run_many_tables_within_a_distance_of_1_from_the_articles() {
	withDatabase(func(a *sqlite3.Adapter) {
		createManyExampleTables(a)

		schema, err := LoadSchema(a)
		if err != nil {
			panic(err)
		}

		generator := ErdGenerator{Table: "articles", Distance: 1}
		generator.Run(schema)

		// Output:
		// entity articles {
		//   * id : integer
		//   --
		//   * user_id : integer
		//   --
		//   index_user_id_on_articles (user_id)
		// }
		//
		// entity comments {
		//   * id : integer
		//   --
		//   * article_id : integer
		//   --
		//   index_article_id_on_articles (article_id)
		// }
		//
		// entity likes {
		//   * article_id : integer
		//   * user_id : integer
		//   --
		//   index_user_id_on_likes (user_id)
		//   - index_article_id_and_user_id_on_likes (article_id, user_id)
		// }
		//
		// entity revisions {
		//   * id : integer
		//   --
		//   * article_id : integer
		//   --
		//   index_article_id_on_revisions (article_id)
		// }
		//
		// entity users {
		//   * id : integer
		//   --
		//   name : text
		// }
		//
		// articles }-- users
		//
		// comments }-- articles
		//
		// likes }-- users
		//
		// likes }-- articles
		//
		// revisions }-- articles
	})
}
