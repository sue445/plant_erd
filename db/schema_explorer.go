package db

import (
	"github.com/deckarep/golang-set/v2"
	"sort"
)

// SchemaExplorer represents schema explorer
type SchemaExplorer struct {
	schema *Schema
	graph  *UndirectedGraph
}

// NewSchemaExplorer returns a new SchemaExplorer instance
func NewSchemaExplorer(schema *Schema) *SchemaExplorer {
	graph := NewUndirectedGraph()
	for _, table := range schema.Tables {
		for _, foreignKey := range table.ForeignKeys {
			graph.PutSymmetric(table.Name, foreignKey.ToTable, true)
		}
	}

	return &SchemaExplorer{schema: schema, graph: graph}
}

// Explore returns surrounding tables from table
func (e *SchemaExplorer) Explore(tableName string, distance int) []string {
	if distance < 0 {
		distance = 0
	}

	foundTableNames := mapset.NewSet[string]()

	e.explore(tableName, distance, foundTableNames, 0)

	tableNames := foundTableNames.ToSlice()

	sort.Strings(tableNames)
	return tableNames
}

func (e *SchemaExplorer) explore(tableName string, distance int, foundTableNames mapset.Set[string], pos int) {
	if pos > distance || foundTableNames.Contains(tableName) {
		return
	}
	foundTableNames.Add(tableName)

	for _, aroundTableName := range e.graph.GetRowColumns(tableName) {
		e.explore(aroundTableName, distance, foundTableNames, pos+1)
	}
}
