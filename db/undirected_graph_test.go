package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUndirectedGraph_PutSymmetric(t *testing.T) {
	g := NewUndirectedGraph()
	g.PutSymmetric("a", "b", true)

	assert.Equal(t, true, g.matrix["a"]["b"])
	assert.Equal(t, true, g.matrix["b"]["a"])
}

func TestUndirectedGraph_GetRowColumns(t *testing.T) {
	g := NewUndirectedGraph()
	g.PutSymmetric("a", "b", true)
	g.PutSymmetric("a", "c", false)
	g.PutSymmetric("e", "a", true)
	g.PutSymmetric("d", "a", true)
	g.PutSymmetric("c", "d", true)

	got1 := g.GetRowColumns("a")
	assert.Equal(t, []string{"b", "d", "e"}, got1)

	got2 := g.GetRowColumns("unknown")
	assert.Empty(t, got2)
}
