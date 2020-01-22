package oracle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_FormatDSN(t *testing.T) {
	type fields struct {
		Username       string
		Password       string
		Host           string
		Port           int
		ServiceName    string
		Loc            string
		Isolation      string
		Questionph     *bool
		PrefetchRows   int
		PrefetchMemory int
		As             string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Only host",
			fields: fields{
				Host:           "localhost",
				Port:           -1,
				PrefetchRows:   -1,
				PrefetchMemory: -1,
			},
			want: "localhost",
		},
		{
			name: "Without params",
			fields: fields{
				Host:           "localhost",
				Username:       "system",
				Password:       "oracle",
				Port:           1521,
				ServiceName:    "xe",
				PrefetchRows:   -1,
				PrefetchMemory: -1,
			},
			want: "system/oracle@localhost:1521/xe",
		},
		{
			name: "With all params",
			fields: fields{
				Host:           "localhost",
				Username:       "system",
				Password:       "oracle",
				Port:           1521,
				ServiceName:    "xe",
				Loc:            "America/New_York",
				Isolation:      "READONLY",
				Questionph:     Bool(true),
				PrefetchRows:   10,
				PrefetchMemory: 20,
				As:             "SYSDBA",
			},
			want: "system/oracle@localhost:1521/xe?loc=America/New_York&isolation=READONLY&questionph=true&prefetch_rows=10&prefetch_memory=20&as=SYSDBA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Username:       tt.fields.Username,
				Password:       tt.fields.Password,
				Host:           tt.fields.Host,
				Port:           tt.fields.Port,
				ServiceName:    tt.fields.ServiceName,
				Loc:            tt.fields.Loc,
				Isolation:      tt.fields.Isolation,
				Questionph:     tt.fields.Questionph,
				PrefetchRows:   tt.fields.PrefetchRows,
				PrefetchMemory: tt.fields.PrefetchMemory,
				As:             tt.fields.As,
			}

			got := c.FormatDSN()
			assert.Equal(t, tt.want, got)

			// TODO: do after
			// _, err := oci8.ParseDSN(got)
			// assert.NoError(t, err)
		})
	}
}
