package main

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/sue445/plant_erd/adapter/oracle"
	"github.com/sue445/plant_erd/cmd"
	"github.com/sue445/plant_erd/lib"
	"github.com/urfave/cli/v3"
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
	generator := lib.NewErdGenerator()
	commonFlags := cmd.CreateCliCommonFlags(generator)

	oracleConfig := oracle.NewConfig()

	command := &cli.Command{
		Name:    "plant_erd-oracle",
		Version: fmt.Sprintf("%s (build. %s)", Version, Revision),
		Usage:   "ERD exporter with PlantUML and Mermaid format (for oracle)",
		Flags: append(
			commonFlags,
			&cli.StringFlag{
				Name:        "user",
				Usage:       "Oracle `USER`",
				Required:    true,
				Destination: &oracleConfig.Username,
			},
			&cli.StringFlag{
				Name:        "password",
				Usage:       "Oracle `PASSWORD`",
				Required:    false,
				Destination: &oracleConfig.Password,
				Sources:     cli.EnvVars("ORACLE_PASSWORD"),
			},
			&cli.StringFlag{
				Name:        "host",
				Usage:       "Oracle `HOST`",
				Required:    false,
				Destination: &oracleConfig.Host,
				Value:       "localhost",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "Oracle `PORT`",
				Required:    false,
				Destination: &oracleConfig.Port,
				Value:       1521,
			},
			&cli.StringFlag{
				Name:        "service",
				Usage:       "Oracle `SERVICE` name",
				Required:    true,
				Destination: &oracleConfig.ServiceName,
			},
		),
		Action: func(_ context.Context, _ *cli.Command) error {
			adapter, closeDatabase, err := oracle.NewAdapter(oracleConfig)

			if err != nil {
				return errors.WithStack(err)
			}

			defer closeDatabase() //nolint:errcheck

			schema, err := lib.LoadSchema(adapter)
			if err != nil {
				return errors.WithStack(err)
			}

			return generator.Run(schema) //nolint:errcheck
		},
	}

	// Sort commands
	sort.Slice(command.Commands, func(i, j int) bool {
		return command.Commands[i].Name < command.Commands[j].Name
	})

	// Sort sub-command flags
	for _, c := range command.Commands {
		sort.Sort(cli.FlagsByName(c.Flags))
	}

	err := command.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
