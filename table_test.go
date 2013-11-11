package gocli

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestAddColumnsNotBeingStrings(t *testing.T) {
	table := NewTable()
	table.Separator = " "
	table.Add(1, 2, "a")
	assert.Contains(t, table.String(), "1 2 a")
}

func TestStringLength(t *testing.T) {
	str := Green("ok")
	assert.Equal(t, stringLength(str), 2)
	assert.Equal(t, stringLength("ok"), 2)
}
