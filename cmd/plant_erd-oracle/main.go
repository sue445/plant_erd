package main

import (
	"fmt"
	"github.com/sue445/plant_erd/adapter/oracle"
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
	app.Name = "plant_erd-oracle"
	app.Usage = "ERD exporter with PlantUML and Mermaid format (for oracle)"

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

	oracleConfig := oracle.NewConfig()
	app.Flags = append(
		commonFlags,
		cli.StringFlag{
			Name:        "user",
			Usage:       "Oracle `USER`",
			Required:    true,
			Destination: &oracleConfig.Username,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "Oracle `PASSWORD`",
			Required:    false,
			Destination: &oracleConfig.Password,
			EnvVar:      "ORACLE_PASSWORD",
		},
		cli.StringFlag{
			Name:        "host",
			Usage:       "Oracle `HOST`",
			Required:    false,
			Destination: &oracleConfig.Host,
			Value:       "localhost",
		},
		cli.IntFlag{
			Name:        "port",
			Usage:       "Oracle `PORT`",
			Required:    false,
			Destination: &oracleConfig.Port,
			Value:       1521,
		},
		cli.StringFlag{
			Name:        "service",
			Usage:       "Oracle `SERVICE` name",
			Required:    true,
			Destination: &oracleConfig.ServiceName,
		},
	)

	app.Action = func(c *cli.Context) error {
		adapter, close, err := oracle.NewAdapter(oracleConfig)

		if err != nil {
			return err
		}

		defer close()

		schema, err := lib.LoadSchema(adapter)
		if err != nil {
			return err
		}

		return generator.Run(schema)
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
