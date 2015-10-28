package main

import (
	"reflect"
	"testing"
)

var justifyTests = []struct {
	justify Justify
	width   int
	src     string
	dst     string
}{
	// normal
	{JustLeft, 5, "abc", "abc  "},
	{JustRight, 5, "abc", "  abc"},
	{JustCenter, 5, "abc", " abc "},
	{JustLeft, 6, "abc", "abc   "},
	{JustRight, 6, "abc", "   abc"},
	{JustCenter, 6, "abc", " abc  "},

	// equal
	{JustLeft, 5, "abcde", "abcde"},
	{JustRight, 5, "abcde", "abcde"},
	{JustCenter, 5, "abcde", "abcde"},
	{JustLeft, 6, "abcdef", "abcdef"},
	{JustRight, 6, "abcdef", "abcdef"},
	{JustCenter, 6, "abcdef", "abcdef"},

	// min
	{JustLeft, 2, "a", "a "},
	{JustRight, 2, "a", " a"},
	{JustCenter, 2, "a", "a "},

	// over
	{JustLeft, 4, "abcdefg", "abcdefg"},
	{JustRight, 4, "abcdefg", "abcdefg"},
	{JustCenter, 4, "abcdefg", "abcdefg"},

	// minus
	{JustLeft, -5, "abc", "abc"},
	{JustRight, -5, "abc", "abc"},
	{JustCenter, -5, "abc", "abc"},

	// double-width character
	{JustLeft, 8, "日本語", "日本語  "},
	{JustRight, 8, "日本語", "  日本語"},
	{JustCenter, 8, "日本語", " 日本語 "},
	{JustLeft, 9, "日本語", "日本語   "},
	{JustRight, 9, "日本語", "   日本語"},
	{JustCenter, 9, "日本語", " 日本語  "},
}

func TestJustfy(t *testing.T) {
	for _, test := range justifyTests {
		expect := test.dst
		actual := test.justify.Just(test.width, test.src)
		if actual != expect {
			kind := "lrc"[test.justify]
			t.Errorf("%q.Justfy(%v, %q) = %q; want %q",
				kind, test.width, test.src, actual, expect)
		}
	}
}

var paddingUpdateWidthTests = []struct {
	a      []string
	before []int
	after  []int
}{
	// update
	{[]string{"aaa", "aaa"},
		[]int{1, 2}, []int{3, 3}},
	{[]string{"aaa", "aaa", "aaaa"},
		[]int{1, 2, 3}, []int{3, 3, 4}},

	// update with stretch
	{[]string{"a", "aa"},
		[]int{}, []int{1, 2}},
	{[]string{"a", "aa", "aaa"},
		[]int{1, 2}, []int{1, 2, 3}},

	// no update
	{nil,
		[]int{1, 1}, []int{1, 1}},
	{[]string{"a", "a"},
		[]int{1, 1}, []int{1, 1}},
	{[]string{"aaaaa", "aa"},
		[]int{3, 3}, []int{5, 3}},
	{[]string{"aaa", "aaa", "aaa"},
		[]int{4, 4, 4}, []int{4, 4, 4}},

	// correspondence double-width character
	{[]string{"日本語", "日本語日本語"},
		[]int{1, 1}, []int{6, 12}},
	{[]string{"「おはよう」", "『こんにちは』"},
		[]int{1, 1}, []int{12, 14}},
}

func TestPaddingUpdateWidth(t *testing.T) {
	for _, test := range paddingUpdateWidthTests {
		p, _ := NewPadding("")
		p.width = test.before
		p.UpdateWidth(test.a)

		expect := test.after
		actual := p.width
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%v.UpdateWidth(%v) got %v; want %v",
				test.before, test.a, actual, expect)
		}
	}
}

var paddingJustKindTests = []struct {
	seq string
	src int
	dst Justify
}{
	{"", -1, JustLeft},
	{"", 0, JustLeft},
	{"", 1, JustLeft},
	{"", 2, JustLeft},

	{"l", 0, JustLeft},
	{"r", 0, JustRight},
	{"c", 0, JustCenter},

	{"r", -2, JustRight},
	{"r", -1, JustRight},
	{"r", 1, JustRight},
	{"r", 2, JustRight},

	{"rc", -1, JustRight},
	{"rc", 0, JustRight},
	{"rc", 1, JustCenter},
	{"rc", 2, JustCenter},
	{"rc", 3, JustCenter},

	{"lcr", -1, JustLeft},
	{"lcr", 0, JustLeft},
	{"lcr", 1, JustCenter},
	{"lcr", 2, JustRight},
	{"lcr", 3, JustCenter},
	{"lcr", 4, JustRight},
	{"lcr", 5, JustCenter},

	{"lcrr", -1, JustLeft},
	{"lcrr", 0, JustLeft},
	{"lcrr", 1, JustCenter},
	{"lcrr", 2, JustRight},
	{"lcrr", 3, JustRight},
	{"lcrr", 4, JustCenter},
	{"lcrr", 5, JustRight},
	{"lcrr", 6, JustRight},
	{"lcrr", 7, JustCenter},
}

func TestsPaddingJustKind(t *testing.T) {
	for _, test := range paddingJustKindTests {
		p, err := NewPadding(test.seq)
		if err != nil {
			t.Errorf("NewPadding(%q) returns %q, want nil",
				test.seq, err)
			continue
		}

		expect := test.dst
		actual := p.justKind(test.src)
		if actual != expect {
			t.Errorf("NewPadding(%q).justkind(%v) = %v; want %v",
				test.seq, test.src, actual, expect)
		}
	}
}

var paddingFormatTests = []struct {
	seq   string
	width []int
	src   []string
	dst   []string
}{
	{"", nil,
		[]string{"a"},
		[]string{"a"}},
	{"", nil,
		[]string{"a", "b"},
		[]string{"a", "b"}},

	{"", []int{2, 3},
		[]string{"a", "b"},
		[]string{"a ", "b  "}},
	{"", []int{2, 3, 4},
		[]string{"a", "b", "c"},
		[]string{"a ", "b  ", "c   "}},
	{"", []int{2, 3, 4, 5},
		[]string{"a", "b", "b", "d"},
		[]string{"a ", "b  ", "b   ", "d    "}},
	{"", []int{2, 3, 4, 5, 6},
		[]string{"a", "b", "b", "d", "e"},
		[]string{"a ", "b  ", "b   ", "d    ", "e     "}},

	{"r", []int{2, 3},
		[]string{"a", "b"},
		[]string{" a", "  b"}},
	{"r", []int{2, 3, 4},
		[]string{"a", "b", "c"},
		[]string{" a", "  b", "   c"}},
	{"r", []int{2, 3, 4, 5},
		[]string{"a", "b", "c", "d"},
		[]string{" a", "  b", "   c", "    d"}},
	{"r", []int{2, 3, 4, 5, 6},
		[]string{"a", "b", "c", "d", "e"},
		[]string{" a", "  b", "   c", "    d", "     e"}},

	{"lcr", []int{5, 1, 5},
		[]string{"n", "=", "100"},
		[]string{"n    ", "=", "  100"}},
	{"lcr", []int{8, 1, 8},
		[]string{"n", "=", "100"},
		[]string{"n       ", "=", "     100"}},
}

func TestPaddingFormat(t *testing.T) {
	for _, test := range paddingFormatTests {
		p, err := NewPadding(test.seq)
		if err != nil {
			t.Errorf("NewPadding(%q) returns %q, want nil",
				test.seq, err)
			continue
		}
		p.width = test.width

		expect := test.dst
		actual := p.Format(test.src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("NewPadding(%q, %v).Format(%q) = %q; want %q",
				test.seq, test.width, test.src, actual, expect)
		}
	}
}
