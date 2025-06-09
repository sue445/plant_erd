package cmd

import (
	"github.com/sue445/plant_erd/lib"
	"github.com/urfave/cli/v3"
)

// CreateCliCommonFlags returns common flags for cli
func CreateCliCommonFlags(generator *lib.ErdGenerator) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "`FILE` for output (default: stdout)",
			Required:    false,
			Destination: &generator.Filepath,
		},
		&cli.StringFlag{
			Name:        "table",
			Aliases:     []string{"t"},
			Usage:       "Output only tables within a certain distance adjacent to each other with foreign keys from a specific `TABLE`",
			Required:    false,
			Destination: &generator.Table,
		},
		&cli.IntFlag{
			Name:        "distance",
			Aliases:     []string{"d"},
			Usage:       "Output only tables within a certain `DISTANCE` adjacent to each other with foreign keys from a specific table",
			Required:    false,
			Destination: &generator.Distance,
			Value:       0,
		},
		&cli.BoolFlag{
			Name:        "skip-index",
			Aliases:     []string{"i"},
			Usage:       "Whether don't print index to ERD. This option is used only --format=plant_uml",
			Required:    false,
			Destination: &generator.SKipIndex,
		},
		&cli.StringFlag{
			Name:        "skip-table",
			Aliases:     []string{"s"},
			Usage:       "Skip generating table by using regex patterns",
			Required:    false,
			Destination: &generator.SkipTable,
		},
		&cli.StringFlag{
			Name:        "format",
			Usage:       "Output format (plant_uml, mermaid. default:plant_uml)",
			Required:    false,
			Destination: &generator.Format,
		},
		&cli.BoolFlag{
			Name:        "show-comment",
			Usage:       "Show column comment. This option is used only --format=mermaid",
			Required:    false,
			Destination: &generator.ShowComment,
		},
	}
}
