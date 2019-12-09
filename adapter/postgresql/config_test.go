package postgresql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_FormatDSN(t *testing.T) {
	type fields struct {
		DBName   string
		User     string
		Password string
		Host     string
		Port     int
		SslMode  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "full fields",
			fields: fields{
				DBName:   "plant_erd_test",
				User:     "postgres",
				Password: "password",
				Host:     "localhost",
				Port:     5432,
				SslMode:  "verify-full",
			},
			want: "dbname=plant_erd_test user=postgres password=password host=localhost port=5432 sslmode=verify-full",
		},
		{
			name:   "no fields",
			fields: fields{},
			want:   "",
		},
		{
			name: "1 field",
			fields: fields{
				DBName: "plant_erd_test",
			},
			want: "dbname=plant_erd_test",
		},
		{
			name: "contains space",
			fields: fields{
				Password: "with spaces",
			},
			want: "password='with spaces'",
		},
		{
			name: "contains space and '",
			fields: fields{
				Password: "it's valid",
			},
			want: "password='it\\'s valid'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				DBName:   tt.fields.DBName,
				User:     tt.fields.User,
				Password: tt.fields.Password,
				Host:     tt.fields.Host,
				Port:     tt.fields.Port,
				SslMode:  tt.fields.SslMode,
			}

			got := c.FormatDSN()
			assert.Equal(t, tt.want, got)
		})
	}
}
