package main

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/sue445/plant_erd/adapter"
	"github.com/sue445/plant_erd/adapter/mysql"
	"github.com/sue445/plant_erd/adapter/postgresql"
	"github.com/sue445/plant_erd/adapter/sqlite3"
	"github.com/sue445/plant_erd/db"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string
)

func main() {
	app := cli.NewApp()

	app.Version = fmt.Sprintf("%s (build. %s)", Version, Revision)
	app.Name = "plant_erd"
	app.Usage = "ERD exporter with PlantUML format"

	generator := ErdGenerator{}

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "f,file",
			Usage:       "Filepath for output (default. stdout)",
			Required:    false,
			Destination: &generator.Filepath,
		},
		cli.StringFlag{
			Name:        "t,table",
			Usage:       "Search surrounding tables",
			Required:    false,
			Destination: &generator.Table,
		},
		cli.IntFlag{
			Name:        "d,distance",
			Usage:       "Search surrounding tables within distance",
			Required:    false,
			Destination: &generator.Distance,
			Value:       0,
		},
	}

	sqlite3Database := ""
	mysqlConfig := mysqlDriver.NewConfig()
	mysqlHost := ""
	mysqlPort := 0
	postgresqlConfig := postgresql.NewConfig()

	app.Commands = []cli.Command{
		{
			Name:    "sqlite3",
			Aliases: []string{"s"},
			Usage:   "Output erd from sqlite3",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "database",
					Usage:       "SQLite3 Database file",
					Required:    true,
					Destination: &sqlite3Database,
				},
			),
			Action: func(c *cli.Context) error {
				adapter, close, err := sqlite3.NewAdapter(sqlite3Database)

				if err != nil {
					return err
				}

				defer close()

				schema, err := loadSchema(adapter)
				if err != nil {
					return err
				}

				return generator.Run(schema)
			},
		},
		{
			Name:    "mysql",
			Aliases: []string{"m"},
			Usage:   "Output erd from mysql",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "host",
					Usage:       "MySQL host",
					Required:    false,
					Destination: &mysqlHost,
					Value:       "localhost",
				},
				cli.IntFlag{
					Name:        "port",
					Usage:       "MySQL port",
					Required:    false,
					Destination: &mysqlPort,
					Value:       3306,
				},
				cli.StringFlag{
					Name:        "user",
					Usage:       "MySQL user",
					Required:    false,
					Destination: &mysqlConfig.User,
					Value:       "root",
				},
				cli.StringFlag{
					Name:        "password",
					Usage:       "MySQL password",
					Required:    false,
					Destination: &mysqlConfig.Passwd,
					EnvVar:      "MYSQL_PASSWORD",
				},
				cli.StringFlag{
					Name:        "database",
					Usage:       "MySQL database name",
					Required:    true,
					Destination: &mysqlConfig.DBName,
				},
				cli.StringFlag{
					Name:        "collation",
					Usage:       "MySQL collation",
					Required:    false,
					Destination: &mysqlConfig.Collation,
					Value:       "utf8_general_ci",
				},
			),
			Action: func(c *cli.Context) error {
				mysqlConfig.Net = "tcp"
				mysqlConfig.Addr = fmt.Sprintf("%s:%d", mysqlHost, mysqlPort)

				adapter, close, err := mysql.NewAdapter(mysqlConfig)

				if err != nil {
					return err
				}

				defer close()

				schema, err := loadSchema(adapter)
				if err != nil {
					return err
				}

				return generator.Run(schema)
			},
		},
		{
			Name:    "postgresql",
			Aliases: []string{"p"},
			Usage:   "Output erd from PostgreSQL",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "host",
					Usage:       "PostgreSQL host",
					Required:    false,
					Destination: &postgresqlConfig.Host,
					Value:       "localhost",
				},
				cli.IntFlag{
					Name:        "port",
					Usage:       "PostgreSQL port",
					Required:    false,
					Destination: &postgresqlConfig.Port,
					Value:       5432,
				},
				cli.StringFlag{
					Name:        "user",
					Usage:       "PostgreSQL user",
					Required:    false,
					Destination: &postgresqlConfig.User,
				},
				cli.StringFlag{
					Name:        "password",
					Usage:       "PostgreSQL password",
					Required:    false,
					Destination: &postgresqlConfig.Password,
					EnvVar:      "POSTGRES_PASSWORD",
				},
				cli.StringFlag{
					Name:        "database",
					Usage:       "PostgreSQL database name",
					Required:    true,
					Destination: &postgresqlConfig.DBName,
				},
				cli.StringFlag{
					Name:        "sslmode",
					Usage:       "PostgreSQL sslmode. c.f. https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS",
					Required:    false,
					Destination: &postgresqlConfig.SslMode,
					Value:       "disable",
				},
			),
			Action: func(c *cli.Context) error {
				adapter, close, err := postgresql.NewAdapter(postgresqlConfig)

				if err != nil {
					return err
				}

				defer close()

				schema, err := loadSchema(adapter)
				if err != nil {
					return err
				}

				return generator.Run(schema)
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	// Sort sub-command flags
	for _, c := range app.Commands {
		sort.Sort(cli.FlagsByName(c.Flags))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadSchema(adapter adapter.Adapter) (*db.Schema, error) {
	tableNames, err := adapter.GetAllTableNames()
	if err != nil {
		return nil, err
	}

	var tables []*db.Table
	for _, tableName := range tableNames {
		table, err := adapter.GetTable(tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return db.NewSchema(tables), nil
}
