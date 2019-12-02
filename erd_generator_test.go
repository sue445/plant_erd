package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/plant_erd/db"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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
					Sequence:   0,
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
