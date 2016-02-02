package gocli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func NewTable() *Table {
	return &Table{
		Columns:   [][]string{},
		Separator: "\t",
	}
}

type Table struct {
	Columns   [][]string
	Separator string

	header []string
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
	toPrint := t.Columns
	if len(t.header) > 0 {
		toPrint = append([][]string{t.header}, toPrint...)
	}

	lengths := map[int]int{}
	for _, c := range toPrint {
		for i, v := range c {
			if l := stringLength(v); l > lengths[i] {
				lengths[i] = l
			}
		}
	}
	for row, col := range toPrint {
		cl := []string{}
		if printIndex {
			col = append([]string{strconv.Itoa(row + 1)}, col...)
		}
		for i, v := range col {
			theLen := lengths[i]
			if printIndex {
				if i == 0 {
					theLen = intLength(len(t.Columns))
				} else {
					theLen = lengths[i-1]
				}
			}
			pad := theLen - stringLength(v)
			cl = append(cl, v+strings.Repeat(" ", pad))
		}
		lines = append(lines, strings.Join(cl, t.Separator))
	}
	return
}

func (t *Table) toPrint() [][]string {
	toPrint := t.Columns
	if len(t.header) > 0 {
		toPrint = append([][]string{t.header}, toPrint...)
	}
	return toPrint
}

func (t *Table) lengths() map[int]int {
	m := map[int]int{}
	for _, c := range t.toPrint() {
		for i, v := range c {
			if l := stringLength(v); l > m[i] {
				m[i] = l
			}
		}
	}
	return m
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
	t.Columns = append(t.Columns, list)
}

// Add adds a column to the table
func (t *Table) Add(cols ...interface{}) {
	t.AddStrings(convertIArray(cols...))
}

func convertIArray(in ...interface{}) (out []string) {
	out = make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, vToS(v))
	}
	return out
}

func vToS(in interface{}) string {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "<nil>"
	}
	switch c := in.(type) {
	case *string:
		return *c
	case *int:
		return fmt.Sprint(*c)
	case *int64:
		return fmt.Sprint(*c)
	case *float64:
		return fmt.Sprint(*c)
	case *time.Time:
		return fmt.Sprint(*c)
	case *bool:
		return fmt.Sprint(*c)
	default:
		return fmt.Sprint(in)
	}
}

func (t *Table) Header(cols ...interface{}) {
	t.header = convertIArray(cols...)
}

// Dereferencing pointers if not nil
// TODO: Please someone tell me, how to do this right!
func (t *Table) AddP(cols ...interface{}) {
	converted := make([]string, 0, len(cols))
	var str string
	for _, v := range cols {
		if value := reflect.ValueOf(v); value.Kind() == reflect.Ptr {
			indirect := reflect.Indirect(value)
			switch {
			case indirect != reflect.Zero(value.Type()) && indirect.IsValid() == true:
				switch {
				case indirect.Kind() == reflect.String:
					str = fmt.Sprint(indirect.String())
				case indirect.Kind() == reflect.Int:
					str = fmt.Sprint(indirect.Int())
				case indirect.Kind() == reflect.Float32:
					str = fmt.Sprint(indirect.Float())
				case indirect.Kind() == reflect.Bool:
					str = fmt.Sprint(indirect.Bool())
				case indirect.Kind() == reflect.Slice:
					str = fmt.Sprint(v)
				default:
					str = fmt.Sprint(v)
				}
			default:
				str = ""
			}
		} else {
			str = fmt.Sprint(v)
		}
		converted = append(converted, str)
	}
	t.AddStrings(converted)
}
