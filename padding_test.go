package alita

import (
	"reflect"
	"testing"
)

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
