package oracle

import (
	"fmt"
	"strconv"
	"strings"
)

// Config represents configuration for Oracle connection
type Config struct {
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

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{Port: -1, PrefetchRows: -1, PrefetchMemory: -1}
}

// FormatDSN formats the given Config into a DSN string which can be passed to the driver.
func (c Config) FormatDSN() string {
	dsn := ""

	if len(c.Username) > 0 {
		dsn += c.Username

		if len(c.Password) > 0 {
			dsn += "/" + c.Password
		}

		dsn += "@"
	}

	dsn += c.Host

	if c.Port >= 0 {
		dsn += ":" + strconv.Itoa(c.Port)
	}

	if len(c.ServiceName) > 0 {
		dsn += "/" + c.ServiceName
	}

	var params []string

	if len(c.Loc) > 0 {
		params = append(params, "loc="+c.Loc)
	}

	if len(c.Isolation) > 0 {
		params = append(params, "isolation="+c.Isolation)
	}

	if c.Questionph != nil {
		params = append(params, fmt.Sprintf("questionph=%v", *c.Questionph))
	}

	if c.PrefetchRows >= 0 {
		params = append(params, "prefetch_rows="+strconv.Itoa(c.PrefetchRows))
	}

	if c.PrefetchMemory >= 0 {
		params = append(params, "prefetch_memory="+strconv.Itoa(c.PrefetchMemory))
	}

	if len(c.As) > 0 {
		params = append(params, "as="+c.As)
	}

	if len(params) > 0 {
		dsn += "?" + strings.Join(params, "&")
	}

	return dsn
}

// Bool returns pointer of bool
func Bool(b bool) *bool {
	return &b
}
