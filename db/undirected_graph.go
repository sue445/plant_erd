package db

import "sort"

// UndirectedGraph represents undirected graph
type UndirectedGraph struct {
	matrix map[string]map[string]bool
}

// NewUndirectedGraph returns a new UndirectedGraph instance
func NewUndirectedGraph() *UndirectedGraph {
	return &UndirectedGraph{matrix: map[string]map[string]bool{}}
}

// PutSymmetric put value to symmetric matrix
func (g *UndirectedGraph) PutSymmetric(row string, col string, value bool) {
	g.initRow(row)
	g.initRow(col)

	g.matrix[row][col] = value
	g.matrix[col][row] = value
}

func (g *UndirectedGraph) initRow(row string) {
	_, ok := g.matrix[row]

	if ok {
		return
	}

	g.matrix[row] = map[string]bool{}
}

// GetRowColumns returns columns of row
func (g *UndirectedGraph) GetRowColumns(row string) []string {
	var columns []string

	g.initRow(row)

	for k, v := range g.matrix[row] {
		if v {
			columns = append(columns, k)
		}
	}

	sort.Strings(columns)
	return columns
}
