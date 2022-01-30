package main

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/sue445/plant_erd/adapter/mysql"
	"github.com/sue445/plant_erd/adapter/postgresql"
	"github.com/sue445/plant_erd/adapter/sqlite3"
	"github.com/sue445/plant_erd/lib"
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
	app.Usage = "ERD exporter with PlantUML and Mermaid format"

	generator := lib.ErdGenerator{ShowComment: true}

	commonFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "f,file",
			Usage:       "`FILE` for output (default: stdout)",
			Required:    false,
			Destination: &generator.Filepath,
		},
		cli.StringFlag{
			Name:        "t,table",
			Usage:       "Output only tables within a certain distance adjacent to each other with foreign keys from a specific `TABLE`",
			Required:    false,
			Destination: &generator.Table,
		},
		cli.IntFlag{
			Name:        "d,distance",
			Usage:       "Output only tables within a certain `DISTANCE` adjacent to each other with foreign keys from a specific table",
			Required:    false,
			Destination: &generator.Distance,
			Value:       0,
		},
		cli.BoolFlag{
			Name:        "i,skip-index",
			Usage:       "Whether don't print index to ERD. This option is used only --format=plant_uml",
			Required:    false,
			Destination: &generator.SKipIndex,
		},
		cli.StringFlag{
			Name:        "s,skip-table",
			Usage:       "Skip generating table by using regex patterns",
			Required:    false,
			Destination: &generator.SkipTable,
		},
		cli.StringFlag{
			Name:        "format",
			Usage:       "Output format (plant_uml, mermaid. default:plant_uml)",
			Required:    false,
			Destination: &generator.Format,
		},
		cli.BoolFlag{
			Name:        "show-comment",
			Usage:       "Show column comment. This option is used only --format=mermaid",
			Required:    false,
			Destination: &generator.ShowComment,
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
			Usage:   "Generate ERD from sqlite3",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "database",
					Usage:       "SQLite3 `DATABASE` file",
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

				schema, err := lib.LoadSchema(adapter)
				if err != nil {
					return err
				}

				return generator.Run(schema)
			},
		},
		{
			Name:    "mysql",
			Aliases: []string{"m"},
			Usage:   "Generate ERD from mysql",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "host",
					Usage:       "MySQL `HOST`",
					Required:    false,
					Destination: &mysqlHost,
					Value:       "localhost",
				},
				cli.IntFlag{
					Name:        "port",
					Usage:       "MySQL `PORT`",
					Required:    false,
					Destination: &mysqlPort,
					Value:       3306,
				},
				cli.StringFlag{
					Name:        "user",
					Usage:       "MySQL `USER`",
					Required:    false,
					Destination: &mysqlConfig.User,
					Value:       "root",
				},
				cli.StringFlag{
					Name:        "password",
					Usage:       "MySQL `PASSWORD`",
					Required:    false,
					Destination: &mysqlConfig.Passwd,
					EnvVar:      "MYSQL_PASSWORD",
				},
				cli.StringFlag{
					Name:        "database",
					Usage:       "MySQL `DATABASE` name",
					Required:    true,
					Destination: &mysqlConfig.DBName,
				},
				cli.StringFlag{
					Name:        "collation",
					Usage:       "MySQL `COLLATION`",
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

				schema, err := lib.LoadSchema(adapter)
				if err != nil {
					return err
				}

				return generator.Run(schema)
			},
		},
		{
			Name:    "postgresql",
			Aliases: []string{"p"},
			Usage:   "Generate ERD from PostgreSQL",
			Flags: append(
				commonFlags,
				cli.StringFlag{
					Name:        "host",
					Usage:       "PostgreSQL `HOST`",
					Required:    false,
					Destination: &postgresqlConfig.Host,
					Value:       "localhost",
				},
				cli.IntFlag{
					Name:        "port",
					Usage:       "PostgreSQL `PORT`",
					Required:    false,
					Destination: &postgresqlConfig.Port,
					Value:       5432,
				},
				cli.StringFlag{
					Name:        "user",
					Usage:       "PostgreSQL `USER`",
					Required:    false,
					Destination: &postgresqlConfig.User,
				},
				cli.StringFlag{
					Name:        "password",
					Usage:       "PostgreSQL `PASSWORD`",
					Required:    false,
					Destination: &postgresqlConfig.Password,
					EnvVar:      "POSTGRES_PASSWORD",
				},
				cli.StringFlag{
					Name:        "database",
					Usage:       "PostgreSQL `DATABASE` name",
					Required:    true,
					Destination: &postgresqlConfig.DBName,
				},
				cli.StringFlag{
					Name:        "sslmode",
					Usage:       "PostgreSQL `SSLMODE`. c.f. https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS",
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

				schema, err := lib.LoadSchema(adapter)
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
