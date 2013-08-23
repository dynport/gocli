package gocli

import (
	"fmt"
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

func (self *Table) String() string {
	lines := []string{}
	for _, line := range self.Lines() {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (self *Table) Lines() (lines []string) {
	for _, col := range self.Columns {
		cl := []string{}
		for i, v := range col {
			cl = append(cl, fmt.Sprintf("%-*s", self.Lengths[i], v))
		}
		lines = append(lines, strings.Join(cl, self.Separator))
	}
	return
}

func (t *Table) AddStrings(s []string) {
	var ret = make([]interface{}, len(s))

	for i, v := range s {
		ret[i] = v
	}
	t.Add(ret)
}

// Add adds a column to the table
func (self *Table) Add(cols ...interface{}) {
	converted := make([]string, len(cols))
	for i, v := range cols {
		s := fmt.Sprint(v)
		converted[i] = s
		if self.Lengths[i] < len(s) {
			self.Lengths[i] = len(s)
		}
	}
	self.Columns = append(self.Columns, converted)
}
