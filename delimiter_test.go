package alita

import (
	"reflect"
	"regexp"
	"testing"
)

type ExprTest struct {
	useRegexp bool
	expr      string
	re        *regexp.Regexp
}

var indexTestsExpr = []ExprTest{
	{false, `=`, regexp.MustCompile(`=`)},
	{false, `=+`, regexp.MustCompile(`=\+`)},
	{false, `-*>`, regexp.MustCompile(`-\*>`)},

	{true, `=`, regexp.MustCompile(`=`)},
	{true, `=+`, regexp.MustCompile(`=+`)},
	{true, `-*>`, regexp.MustCompile(`-*>`)},
}

func TestParseExpr(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsExpr {
		d.UseRegexp = test.useRegexp
		if err := d.Set(test.expr); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.expr, err)
		}
		if !reflect.DeepEqual(d.re, test.re) {
			t.Errorf("got %q, want %q",
				d.re, test.re)
		}
	}
}

type DefaultSplitTest struct {
	src string
	dst []string
}

var indexTestsDefaultSplit = []DefaultSplitTest{
	{"a", []string{"a"}},
	{"a b", []string{"a", "b"}},
	{"a b c", []string{"a", "b", "c"}},
	{"a b c d", []string{"a", "b", "c", "d"}},

	{"a  b", []string{"a", "b"}},
	{"a  b c", []string{"a", "b", "c"}},
}

func TestDefaultSplit(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsDefaultSplit {
		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q).Split(%q) = %q; want %q",
				d.re, test.src, actual, expect)
		}
	}
}

type SplitTest struct {
	useRegexp bool
	expr      string
	src       string
	dst       []string
}

var indexTestsSplit = []SplitTest{
	{false, `=`, "n", []string{"n"}},
	{false, `=`, "n=", []string{"n", "=", ""}},
	{false, `=`, "=n", []string{"", "=", "n"}},

	{false, `=`, "n=100", []string{"n", "=", "100"}},
	{false, `=`, "n  =  100", []string{"n", "=", "100"}},
	{false, `=`, "n  = 100=200", []string{"n", "=", "100", "=", "200"}},
	{false, `=`, "n=100 =  200", []string{"n", "=", "100", "=", "200"}},

	{false, `=`, "n==100", []string{"n", "=", "", "=", "100"}},
	{false, `=`, "n===100", []string{"n", "=", "", "=", "", "=", "100"}},

	{false, `=>`, "a=>b=>c", []string{"a", "=>", "b", "=>", "c"}},
	{false, `=>`, "a => b => c", []string{"a", "=>", "b", "=>", "c"}},

	{true, `=+>`, "a => b",
		[]string{"a", "=>", "b"}},
	{true, `=+>`, "a => b ==> c ===> d",
		[]string{"a", "=>", "b", "==>", "c", "===>", "d"}},
}

func TestSplit(t *testing.T) {
	d := NewDelimiter()
	for _, test := range indexTestsSplit {
		d.UseRegexp = test.useRegexp
		if err := d.Set(test.expr); err != nil {
			t.Errorf("Set(%q) returns %q; want nil",
				test.expr, err)
		}
		actual := d.Split(test.src)
		expect := test.dst
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Delimiter(%q).Split(%q) = %q; want %q",
				d.re, test.src, actual, expect)
		}
	}
}
