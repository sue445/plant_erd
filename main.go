package main

import (
	"fmt"
	"github.com/sue445/plant_erd/adapter"
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
