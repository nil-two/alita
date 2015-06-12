package main

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

var marginSetTests = []struct {
	format string
	left   int
	right  int
}{
	// digit only
	{"1", 1, 1},
	{"2", 2, 2},
	{"10", 10, 10},
	{"00", 0, 0},
	{"001", 1, 1},

	// colon separated digits
	{"0:3", 0, 3},
	{"00:001", 0, 1},
	{"1:1", 1, 1},
	{"2:1", 2, 1},
	{"1:4", 1, 4},
	{"10:5", 10, 5},
	{"5:20", 5, 20},
}

func TestMarginSet(t *testing.T) {
	m := NewMargin()
	for _, test := range marginSetTests {
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

var marginSetErrTests = []string{
	"abc",
	"100000000000000000000000000000",
	"-1",
	":",
	"1:",
	":1",
	"1:-1",
	"-1:1",
	"::",
	"1::2",
	"3:1:4",
}

func TestMarginSetErr(t *testing.T) {
	m := NewMargin()
	for _, format := range marginSetErrTests {
		if err := m.Set(format); err == nil {
			t.Errorf("Margin.Set(%q) returns nil; want err", format)
		}
	}
}

var marginJoinTests = []struct {
	left  int
	right int
	src   []string
	dst   string
}{
	{0, 0, nil, ""},
	{1, 1, nil, ""},

	{0, -2, []string{"a", "b", "c"}, "abc"},
	{-4, 0, []string{"a", "b", "c"}, "abc"},
	{-4, -5, []string{"a", "b", "c"}, "abc"},

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
	for _, test := range marginJoinTests {
		m.SetMargin(test.left, test.right)

		actual := m.Join(test.src)
		expect := test.dst
		if actual != expect {
			t.Errorf("Margin(%d, %d).Join(%q) = %q; want %q",
				test.left, test.right, test.src, actual, expect)
		}
	}
}
