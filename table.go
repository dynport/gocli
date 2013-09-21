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

func (self *Table) String() string {
	lines := []string{}
	for _, line := range self.Lines() {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

var uncolorRegexp = regexp.MustCompile("\033\\[38;5;\\d+m")

func stringLength(s string) int {
	return len(strings.Replace(uncolorRegexp.ReplaceAllString(s, ""), "\033[0m", "", -1))
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
	t.Add(ret...)
}

// Add adds a column to the table
func (self *Table) Add(cols ...interface{}) {
	converted := make([]string, len(cols))
	for i, v := range cols {
		s := fmt.Sprint(v)
		converted[i] = s
		if self.Lengths[i] < stringLength(s) {
			self.Lengths[i] = len(s)
		}
	}
	self.Columns = append(self.Columns, converted)
}
