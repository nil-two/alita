package alita

import (
	"reflect"
	"testing"
)

type JustfyTest struct {
	justfy Justify
	width  int
	src    string
	dst    string
}

var indexTestsJustfy = []JustfyTest{
	// Case normal
	{JustLeft, 5, "abc", "abc  "},
	{JustRight, 5, "abc", "  abc"},
	{JustCenter, 5, "abc", " abc "},
	{JustLeft, 6, "abc", "abc   "},
	{JustRight, 6, "abc", "   abc"},
	{JustCenter, 6, "abc", " abc  "},

	// Case equal
	{JustLeft, 5, "abcde", "abcde"},
	{JustRight, 5, "abcde", "abcde"},
	{JustCenter, 5, "abcde", "abcde"},
	{JustLeft, 6, "abcdef", "abcdef"},
	{JustRight, 6, "abcdef", "abcdef"},
	{JustCenter, 6, "abcdef", "abcdef"},

	// Case min
	{JustLeft, 2, "a", "a "},
	{JustRight, 2, "a", " a"},
	{JustCenter, 2, "a", "a "},

	// Case over
	{JustLeft, 4, "abcdefg", "abcdefg"},
	{JustRight, 4, "abcdefg", "abcdefg"},
	{JustCenter, 4, "abcdefg", "abcdefg"},

	// Case minus
	{JustLeft, -5, "abc", "abc"},
	{JustRight, -5, "abc", "abc"},
	{JustCenter, -5, "abc", "abc"},

	// Case double-width character
	{JustLeft, 8, "日本語", "日本語  "},
	{JustRight, 8, "日本語", "  日本語"},
	{JustCenter, 8, "日本語", " 日本語 "},
	{JustLeft, 9, "日本語", "日本語   "},
	{JustRight, 9, "日本語", "   日本語"},
	{JustCenter, 9, "日本語", " 日本語  "},
}

func TestJustfy(t *testing.T) {
	for _, test := range indexTestsJustfy {
		actual := test.justfy.Just(test.width, test.src)
		expect := test.dst
		if actual != expect {
			kind := "lrc"[test.justfy]
			t.Errorf("%q.Justfy(%q) = %q; want %q",
				kind, test.src, actual, expect)
		}
	}
}

type UpdateWidthTest struct {
	a      []string
	before []int
	after  []int
}

var indexTestsUpdateWidth = []UpdateWidthTest{
	// Through update
	{[]string{"a"},
		[]int{}, []int{}},
	{[]string{"aaaa"},
		[]int{}, []int{}},
	{[]string{"aaaa"},
		[]int{3, 5}, []int{3, 5}},

	// Update
	{[]string{"aaa", "aaa"},
		[]int{1, 2}, []int{3, 3}},
	{[]string{"aaa", "aaa", "aaaa"},
		[]int{1, 2, 3}, []int{3, 3, 4}},

	// Update with stretch
	{[]string{"a", "aa"},
		[]int{}, []int{1, 2}},
	{[]string{"a", "aa", "aaa"},
		[]int{1, 2}, []int{1, 2, 3}},

	// No update
	{[]string{"a", "a"},
		[]int{1, 1}, []int{1, 1}},
	{[]string{"aaaaa", "aa"},
		[]int{3, 3}, []int{5, 3}},
	{[]string{"aaa", "aaa", "aaa"},
		[]int{4, 4, 4}, []int{4, 4, 4}},

	// Correspondence double-width character
	{[]string{"日本語", "日本語日本語"},
		[]int{1, 1}, []int{6, 12}},
	{[]string{"「おはよう」", "『こんにちは』"},
		[]int{1, 1}, []int{12, 14}},
}

func TestUpdateWidth(t *testing.T) {
	p := NewPadding()
	for _, test := range indexTestsUpdateWidth {
		p.width = test.before
		p.UpdateWidth(test.a)
		actual := p.width
		expect := test.after
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%v -> UpdateWidth(%v) got %v; want %v",
				test.before, test.a, actual, expect)
		}
	}
}

type JustKindTest struct {
	justfies []Justify
	src      int
	dst      Justify
}

var indexTestsJustKind = []JustKindTest{
	{nil, -1, JustLeft},
	{nil, 0, JustLeft},
	{nil, 1, JustLeft},
	{nil, 2, JustLeft},

	{[]Justify{JustLeft}, 0, JustLeft},
	{[]Justify{JustRight}, 0, JustRight},
	{[]Justify{JustCenter}, 0, JustCenter},

	{[]Justify{JustRight}, -2, JustRight},
	{[]Justify{JustRight}, -1, JustRight},
	{[]Justify{JustRight}, 1, JustRight},
	{[]Justify{JustRight}, 2, JustRight},

	{[]Justify{JustRight, JustCenter}, -1, JustRight},
	{[]Justify{JustRight, JustCenter}, 0, JustRight},
	{[]Justify{JustRight, JustCenter}, 1, JustCenter},
	{[]Justify{JustRight, JustCenter}, 2, JustCenter},
	{[]Justify{JustRight, JustCenter}, 3, JustCenter},

	{[]Justify{JustLeft, JustCenter, JustRight}, -1, JustLeft},
	{[]Justify{JustLeft, JustCenter, JustRight}, 0, JustLeft},
	{[]Justify{JustLeft, JustCenter, JustRight}, 1, JustCenter},
	{[]Justify{JustLeft, JustCenter, JustRight}, 2, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight}, 3, JustCenter},
	{[]Justify{JustLeft, JustCenter, JustRight}, 4, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight}, 5, JustCenter},

	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, -1, JustLeft},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 0, JustLeft},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 1, JustCenter},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 2, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 3, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 4, JustCenter},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 5, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 6, JustRight},
	{[]Justify{JustLeft, JustCenter, JustRight, JustRight}, 7, JustCenter},
}

func TestJustKind(t *testing.T) {
	p := NewPadding()
	for _, test := range indexTestsJustKind {
		p.SetJustfies(test.justfies)
		actual := p.justKind(test.src)
		expect := test.dst
		if actual != expect {
			t.Errorf("Padding(%v).justkind(%v) = %v; want %v",
				test.justfies, test.src, actual, expect)
		}
	}
}
