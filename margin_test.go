package alita

import (
	"testing"
)

func TestMarginDefault(t *testing.T) {
	m := NewMargin()
	l, r := 1, 1
	if m.left != l || m.right != r {
		t.Errorf("got %d:%d; want %d:%d", m.left, m.right, l, r)
	}
}

type MarginSetTest struct {
	format string
	left   int
	right  int
}

var indexTestsMarginSet = []MarginSetTest{
	// digit only
	{"1", 1, 1},
	{"2", 2, 2},
	{"10", 10, 10},
	{"00", 0, 0},

	// colon separated digits
	{"1:1", 1, 1},
	{"2:1", 2, 1},
	{"1:4", 1, 4},
	{"10:5", 10, 5},
	{"5:20", 5, 20},
}

func TestMarginSet(t *testing.T) {
	m := NewMargin()
	for _, test := range indexTestsMarginSet {
		if err := m.Set(test.format); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.format, err)
		}
		if m.left != test.left || m.right != test.right {
			t.Errorf("got %d:%d; want %d:%d",
				m.left, m.right, test.left, test.right)
		}
	}
}

var indexTestsMarginSetErr = []string{
	"abc",
	"100000000000000000000000000000",
	"-1",
	":",
	"1:",
	":1",
	"1:-1",
	"-1:1",
	"3:1:4",
}

func TestMarginSetErr(t *testing.T) {
	m := NewMargin()
	for _, format := range indexTestsMarginSetErr {
		if err := m.Set(format); err == nil {
			t.Errorf("Margin.Set(%q) returns nil; want err", format)
		}
	}
}

type MarginJoinTest struct {
	left  int
	right int
	src   []string
	dst   string
}

var indexTestsMarginJoin = []MarginJoinTest{
	{0, 0, []string{"a"}, "a"},
	{2, 2, []string{"a"}, "a"},

	{0, 0, []string{"a", "b"}, "ab"},
	{2, 2, []string{"a", "b"}, "a  b"},

	{1, 1, []string{"n", "=", "100"}, "n = 100"},
	{2, 2, []string{"n", "=", "100"}, "n  =  100"},
	{1, 0, []string{"n", "=", "100"}, "n =100"},
	{0, 1, []string{"n", "=", "100"}, "n= 100"},

	{1, 1, []string{"1", "2", "3", "4"}, "1 2 3 4"},
	{2, 2, []string{"1", "2", "3", "4"}, "1  2  3  4"},
	{0, 1, []string{"1", "2", "3", "4"}, "12 34"},
	{1, 0, []string{"1", "2", "3", "4"}, "1 23 4"},

	{1, 1, []string{"a", ":", "b", ":", "c"}, "a : b : c"},
	{2, 2, []string{"a", ":", "b", ":", "c"}, "a  :  b  :  c"},
	{0, 1, []string{"a", ":", "b", ":", "c"}, "a: b: c"},
	{1, 0, []string{"a", ":", "b", ":", "c"}, "a :b :c"},

	{0, 1, []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		"12 34 56 78"},
	{1, 0, []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		"1 23 45 67 8"},
}

func TestMarginJoin(t *testing.T) {
	m := NewMargin()
	for _, test := range indexTestsMarginJoin {
		m.SetMargin(test.left, test.right)
		actual := m.Join(test.src)
		expect := test.dst
		if actual != expect {
			t.Errorf("Margin(%d, %d).Join(%q) = %q; want %q",
				test.left, test.right, test.src, actual, expect)
		}
	}
}
