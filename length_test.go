package gocli

import "testing"

func TestLength(t *testing.T) {
	tab := NewTable()
	tab.Add("a", "ab")

	lengths := tab.lengths()

	tests := []struct {
		Name     string
		Expected interface{}
		Value    interface{}
	}{
		{"0", 1, lengths[0]},
		{"1", 2, lengths[1]},
	}

	for _, tst := range tests {
		if tst.Expected != tst.Value {
			t.Errorf("expected %s to be %#v, was %#v", tst.Name, tst.Expected, tst.Value)
		}
	}

	tab.Add(Red("abc"), Green("abcd"))

	lengths = tab.lengths()

	tests = []struct {
		Name     string
		Expected interface{}
		Value    interface{}
	}{
		{"0", 3, lengths[0]},
		{"1", 4, lengths[1]},
	}

	for _, tst := range tests {
		if tst.Expected != tst.Value {
			t.Errorf("expected %s to be %#v, was %#v", tst.Name, tst.Expected, tst.Value)
		}
	}
}
