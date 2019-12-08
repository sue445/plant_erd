package postgresql

import (
	"fmt"
	"strconv"
	"strings"
)

// Config represents configuration for PostgreSQL connection
type Config struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     int
	SslMode  string
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{Port: 5432, SslMode: "disable"}
}

// FormatDSN formats the given Config into a DSN string which can be passed to the driver.
func (c Config) FormatDSN() string {
	var params []string

	if len(c.DBName) > 0 {
		params = append(params, "dbname="+c.escape(c.DBName))
	}

	if len(c.User) > 0 {
		params = append(params, "user="+c.escape(c.User))
	}

	if len(c.Password) > 0 {
		params = append(params, "password="+c.escape(c.Password))
	}

	if len(c.Host) > 0 {
		params = append(params, "host="+c.escape(c.Host))
	}

	if c.Port > 0 {
		params = append(params, "port="+strconv.Itoa(c.Port))
	}

	if len(c.SslMode) > 0 {
		params = append(params, "sslmode="+c.escape(c.SslMode))
	}

	return strings.Join(params, " ")
}

func (c *Config) escape(str string) string {
	if !strings.Contains(str, " ") {
		return str
	}

	str = strings.ReplaceAll(str, "'", "\\'")
	return fmt.Sprintf("'%s'", str)
}
