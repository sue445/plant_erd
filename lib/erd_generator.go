package lib

import (
	"fmt"
	"github.com/sue445/plant_erd/db"
	"io/ioutil"
	"os"
	"regexp"
)

// ErdGenerator represents ERD generator
type ErdGenerator struct {
	Filepath    string
	Table       string
	Distance    int
	SKipIndex   bool
	SkipTable   string
	Format      string
	ShowComment bool
}

// NewErdGenerator returns a new NewErdGenerator instance
func NewErdGenerator() *ErdGenerator {
	g := ErdGenerator{ShowComment: true}
	return &g
}

// Run performs generator
func (g *ErdGenerator) Run(schema *db.Schema) error {
	err := g.checkParamTable(schema)

	if err != nil {
		return err
	}

	erd, err := g.generate(schema)
	if err != nil {
		return err
	}

	return g.output(erd)
}

func (g *ErdGenerator) checkParamTable(schema *db.Schema) error {
	if g.Table == "" {
		return nil
	}

	for _, table := range schema.Tables {
		if table.Name == g.Table {
			return nil
		}
	}

	return fmt.Errorf("%s is not found in database", g.Table)
}

func (g *ErdGenerator) generate(schema *db.Schema) (string, error) {
	if g.SkipTable != "" {
		schema = g.filterSchema(schema, []string{g.SkipTable})
	}

	switch g.Format {
	case "", "plant_uml":
		return g.generateErd(schema), nil
	case "mermaid":
		return g.generateMermaid(schema), nil
	}

	return "", fmt.Errorf("%s is unknown format", g.Format)
}

func (g *ErdGenerator) generateErd(schema *db.Schema) string {
	if g.Table == "" || g.Distance <= 0 {
		return schema.ToErd(!g.SKipIndex)
	}

	subset := schema.Subset(g.Table, g.Distance)
	return subset.ToErd(!g.SKipIndex)
}

func (g *ErdGenerator) generateMermaid(schema *db.Schema) string {
	if g.Table == "" || g.Distance <= 0 {
		return schema.ToMermaid(g.ShowComment)
	}

	subset := schema.Subset(g.Table, g.Distance)
	return subset.ToMermaid(g.ShowComment)
}

func (g *ErdGenerator) output(content string) error {
	if g.Filepath == "" {
		// Print to stdout
		fmt.Fprint(os.Stdout, content)
		return nil
	}

	// Output to file
	return ioutil.WriteFile(g.Filepath, []byte(content), 0644)
}

func (g *ErdGenerator) filterSchema(schema *db.Schema, skipTable []string) *db.Schema {
	tableNames := schema.Tables
	var tables []*db.Table
	for _, table := range tableNames {
		if matched := g.matchedSkippedTable(skipTable, table.Name); matched {
			continue
		}
		tables = append(tables, table)
	}
	return db.NewSchema(tables)
}

func (g *ErdGenerator) matchedSkippedTable(skipPatterns []string, tableName string) bool {
	for _, pattern := range skipPatterns {
		if matched, _ := regexp.MatchString(pattern, tableName); matched {
			return true
		}
	}
	return false
}
