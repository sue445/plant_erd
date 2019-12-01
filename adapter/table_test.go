package adapter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable_ToErd(t *testing.T) {
	type fields struct {
		Name        string
		Columns     []*Column
		ForeignKeys []*ForeignKey
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without primary key",
			fields: fields{
				Name: "articles",
				Columns: []*Column{
					{
						Name:    "id",
						Type:    "integer",
						NotNull: true,
					},
					{
						Name:    "user_id",
						Type:    "integer",
						NotNull: true,
					},
					{
						Name: "title",
						Type: "text",
					},
				},
			},
			want: `entity articles {
  * id : integer
  * user_id : integer
  title : text
}`,
		},
		{
			name: "with primary key",
			fields: fields{
				Name: "articles",
				Columns: []*Column{
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
						Name: "title",
						Type: "text",
					},
				},
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
  title : text
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:        tt.fields.Name,
				Columns:     tt.fields.Columns,
				ForeignKeys: tt.fields.ForeignKeys,
			}

			got := table.ToErd()
			assert.Equal(t, tt.want, got)
		})
	}
}
