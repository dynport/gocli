package gocli

import (
	"fmt"
	"strings"
)

func NewTable() *Table {
	return &Table{
		Columns: [][]string{},
		Lengths: map[int]int{},
		Separator: "\t",
	}
}

type Table struct {
	Columns [][]string
	Lengths map[int]int
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

func (self *Table) Add(cols ...string) {
	for i, v := range cols {
		if self.Lengths[i] < len(v) {
			self.Lengths[i] = len(v)
		}
	}
	self.Columns = append(self.Columns, cols)
}
