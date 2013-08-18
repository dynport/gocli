package gocli

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	table := NewTable()
	assert.NotNil(t, table)
	table.Add("a", "b")
	table.Add("aa", "bb")
	assert.Equal(t, len(table.Lines()), 2)
	assert.Contains(t, table.String(), "a\tb")

	table.Separator = " "
	assert.Contains(t, table.String(), "a b")
}

func TestColorize(t *testing.T) {
	str := Colorize(90, "test")
	assert.NotNil(t, str)
}
