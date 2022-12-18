package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_ToErd(t *testing.T) {
	type fields struct {
		Name        string
		Columns     []*Column
		ForeignKeys []*ForeignKey
		Indexes     []*Index
	}
	type args struct {
		showIndex bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
			args: args{
				showIndex: true,
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
			args: args{
				showIndex: true,
			},
			want: `entity articles {
  * id : integer
  --
  * user_id : integer
  title : text
}`,
		},
		{
			name: "with index and enabled showIndex",
			fields: fields{
				Name: "followers",
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
						Name:    "target_user_id",
						Type:    "integer",
						NotNull: true,
					},
				},
				ForeignKeys: []*ForeignKey{
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
				Indexes: []*Index{
					{
						Name:    "index_target_user_id_and_user_id_on_followers",
						Columns: []string{"target_user_id", "user_id"},
						Unique:  true,
					},
					{
						Name:    "index_user_id_and_target_user_id_on_followers",
						Columns: []string{"user_id", "target_user_id"},
						Unique:  false,
					},
				},
			},
			args: args{
				showIndex: true,
			},
			want: `entity followers {
  * id : integer
  --
  * user_id : integer
  * target_user_id : integer
  --
  - index_target_user_id_and_user_id_on_followers (target_user_id, user_id)
  index_user_id_and_target_user_id_on_followers (user_id, target_user_id)
}`,
		},
		{
			name: "with index and disabled showIndex",
			fields: fields{
				Name: "followers",
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
						Name:    "target_user_id",
						Type:    "integer",
						NotNull: true,
					},
				},
				ForeignKeys: []*ForeignKey{
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
				Indexes: []*Index{
					{
						Name:    "index_target_user_id_and_user_id_on_followers",
						Columns: []string{"target_user_id", "user_id"},
						Unique:  true,
					},
					{
						Name:    "index_user_id_and_target_user_id_on_followers",
						Columns: []string{"user_id", "target_user_id"},
						Unique:  false,
					},
				},
			},
			args: args{
				showIndex: false,
			},
			want: `entity followers {
  * id : integer
  --
  * user_id : integer
  * target_user_id : integer
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:        tt.fields.Name,
				Columns:     tt.fields.Columns,
				ForeignKeys: tt.fields.ForeignKeys,
				Indexes:     tt.fields.Indexes,
			}

			got := table.ToErd(tt.args.showIndex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTable_ToMermaid(t *testing.T) {
	type fields struct {
		Name        string
		Columns     []*Column
		ForeignKeys []*ForeignKey
		Indexes     []*Index
	}
	type args struct {
		showComment bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "with comment",
			fields: fields{
				Name: "articles",
				Columns: []*Column{
					{
						Name:       "id",
						Type:       "integer(10) unsigned",
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
				ForeignKeys: []*ForeignKey{
					{
						FromColumn: "user_id",
						ToTable:    "users",
						ToColumn:   "id",
					},
				},
			},
			args: args{
				showComment: true,
			},
			want: `articles {
  integer_10_unsigned id PK "not null"
  integer user_id FK "not null"
  text title
}`,
		},
		{
			name: "without comment",
			fields: fields{
				Name: "articles",
				Columns: []*Column{
					{
						Name:       "id",
						Type:       "integer(10) unsigned",
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
				ForeignKeys: []*ForeignKey{
					{
						FromColumn: "user_id",
						ToTable:    "users",
						ToColumn:   "id",
					},
				},
			},
			args: args{
				showComment: false,
			},
			want: `articles {
  integer_10_unsigned id
  integer user_id
  text title
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &Table{
				Name:        tt.fields.Name,
				Columns:     tt.fields.Columns,
				ForeignKeys: tt.fields.ForeignKeys,
				Indexes:     tt.fields.Indexes,
			}

			got := table.ToMermaid(tt.args.showComment)
			assert.Equal(t, tt.want, got)
		})
	}
}
