package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn_ToErd(t *testing.T) {
	type fields struct {
		Name       string
		Type       string
		NotNull    bool
		PrimaryKey bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without NotNull",
			fields: fields{
				Name: "id",
				Type: "integer",
			},
			want: "id : integer",
		},
		{
			name: "with NotNull",
			fields: fields{
				Name:    "id",
				Type:    "integer",
				NotNull: true,
			},
			want: "* id : integer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Column{
				Name:       tt.fields.Name,
				Type:       tt.fields.Type,
				NotNull:    tt.fields.NotNull,
				PrimaryKey: tt.fields.PrimaryKey,
			}

			got := c.ToErd()
			assert.Equal(t, tt.want, got)
		})
	}
}
