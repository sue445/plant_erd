package main

import (
	"fmt"
	"github.com/sue445/plant_erd/adapter/oracle"
	"github.com/sue445/plant_erd/cmd"
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
	commonFlags := cmd.CreateCliCommonFlags(&generator)

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
