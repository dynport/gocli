package gocli

import (
	"fmt"
	"regexp"
	"strings"
)

func NewTable() *Table {
	return &Table{
		Columns:   [][]string{},
		Lengths:   map[int]int{},
		Separator: "\t",
	}
}

type Table struct {
	Columns   [][]string
	Lengths   map[int]int
	Separator string
}

func (t *Table) String() string {
	return strings.Join(t.Lines(), "\n")
}

var uncolorRegexp = regexp.MustCompile("\033\\[38;5;\\d+m([^\033]+)\033\\[0m")

func stringLength(s string) int {
	return len(uncolorRegexp.ReplaceAllString(s, "$1"))
}

func (t *Table) Lines() (lines []string) {
	for _, col := range t.Columns {
		cl := []string{}
		for i, v := range col {
			cl = append(cl, fmt.Sprintf("%-*s", t.Lengths[i], v))
		}
		lines = append(lines, strings.Join(cl, t.Separator))
	}
	return
}

func (t *Table) AddStrings(list []string) {
	for i, s := range list {
		length := stringLength(s)
		if width := t.Lengths[i]; width < length {
			t.Lengths[i] = length
		}
	}
	t.Columns = append(t.Columns, list)
}

// Add adds a column to the table
func (t *Table) Add(cols ...interface{}) {
	converted := make([]string, 0, len(cols))
	for _, v := range cols {
		converted = append(converted, fmt.Sprint(v))
	}
	t.AddStrings(converted)
}
