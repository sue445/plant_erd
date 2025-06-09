package cmd

import (
	"github.com/sue445/plant_erd/lib"
	"github.com/urfave/cli/v2"
)

// CreateCliCommonFlags returns common flags for cli
func CreateCliCommonFlags(generator *lib.ErdGenerator) []cli.Flag {
	return []cli.Flag{
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
}
