package gocli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
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

	SortBy int
}

func (t *Table) Select(message string) int {
	for {
		fmt.Fprintf(os.Stdout, t.StringWithIndex()+"\n"+message+": ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		i, e := strconv.Atoi(scanner.Text())
		if e == nil {
			if i > 0 && i <= len(t.Columns) {
				return i - 1
			}
		}
	}
}

func (t *Table) Len() int { return len(t.Columns) }

func (t *Table) Swap(a, b int) { t.Columns[a], t.Columns[b] = t.Columns[b], t.Columns[a] }

func (t *Table) Less(a, b int) bool {
	if len(t.Columns[a]) <= t.SortBy {
		return false
	} else if len(t.Columns[b]) <= t.SortBy {
		return true
	} else {
		return fmt.Sprint(t.Columns[a][t.SortBy]) <= fmt.Sprintf(t.Columns[b][t.SortBy])
	}
	return true
}

func (t *Table) String() string {
	return strings.Join(t.Lines(false), "\n")
}

func (t *Table) StringWithIndex() string {
	return strings.Join(t.Lines(true), "\n")
}

var uncolorRegexp = regexp.MustCompile("\033\\[38;5;\\d+m([^\033]+)\033\\[0m")

func stringLength(s string) int {
	return utf8.RuneCountInString((uncolorRegexp.ReplaceAllString(s, "$1")))
}

func (t *Table) Lines(printIndex bool) (lines []string) {
	for row, col := range t.Columns {
		cl := []string{}
		if printIndex {
			col = append([]string{strconv.Itoa(row + 1)}, col...)
		}
		for i, v := range col {
			theLen := t.Lengths[i]
			if printIndex {
				if i == 0 {
					theLen = intLength(len(t.Columns))
				} else {
					theLen = t.Lengths[i-1]
				}
			}
			pad := theLen - stringLength(v)
			cl = append(cl, v+strings.Repeat(" ", pad))
		}
		lines = append(lines, strings.Join(cl, t.Separator))
	}
	return
}

func intLength(i int) int {
	if i == 0 {
		return 1
	} else if i < 0 {
		return intLength(int(math.Abs(float64(i)))) + 1
	}
	return int(math.Ceil(math.Log10(float64(i + 1))))
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
