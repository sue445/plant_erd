package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchema_ToErd(t *testing.T) {
	type fields struct {
		Tables []*Table
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "users and articles",
			fields: fields{
				Tables: []*Table{
					{
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
						},
						ForeignKeys: []*ForeignKey{
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
						Columns: []*Column{
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
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
}

entity users {
  * id : integer
  --
  name : text
}

articles }-- users`,
		},
		{
			name: "Reject foreign key which table isn't in schema",
			fields: fields{
				Tables: []*Table{
					{
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
						},
						ForeignKeys: []*ForeignKey{
							{
								Sequence:   0,
								FromColumn: "user_id",
								ToTable:    "users",
								ToColumn:   "id",
							},
						},
					},
				},
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{
				Tables: tt.fields.Tables,
			}

			got := s.ToErd()
			assert.Equal(t, tt.want, got)
		})
	}
}
